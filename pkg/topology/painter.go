// topology package should do refactoring.
package topology

import (
	"bytes"
	"io/ioutil"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/kwseo/conf-dep-graph/pkg/util"
)

const (
	FILESYSTEM PathType = iota
	URL
)

type GraphvizFunc func(g *graphviz.Graphviz, graph *cgraph.Graph) error

type PathType int

type Config struct {
	Layout   string     `yaml:"layout"`
	Services []*Service `yaml:"services"`
}

func (c *Config) GetLayout() graphviz.Layout {
	switch c.Layout {
	case "circo":
		return graphviz.CIRCO
	case "neato":
		return graphviz.NEATO
	case "dot":
		return graphviz.DOT
	case "fdp":
		return graphviz.FDP
	case "osage":
		return graphviz.OSAGE
	case "patchwork":
		return graphviz.PATCHWORK
	case "sfdp":
		return graphviz.SFDP
	case "twopi":
		return graphviz.TWOPI
	default:
		return graphviz.NEATO
	}
}

func LoadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	config := new(Config)
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

type Topology struct {
	cfg      *Config
	services map[string]*Service
	logger   log.Logger
}

func New(cfg *Config) *Topology {
	return &Topology{
		cfg:      cfg,
		services: map[string]*Service{},
	}
}

func (t *Topology) AddService(newSvc *Service) {
	if _, exist := t.services[newSvc.Name]; exist {
		level.Warn(t.logger).Log("msg", "Already existed service")
		return
	}
	for _, svc := range t.services {
		if svc.DependOn(newSvc) {
			svc.AddDependency(newSvc)
		}
		if newSvc.DependOn(svc) {
			newSvc.AddDependency(svc)
		}
	}
	t.services[newSvc.Name] = newSvc
}

func (t *Topology) GraphAsPNG() ([]byte, error) {
	return t.withGraph(graphviz.PNG, func(g *graphviz.Graphviz, graph *cgraph.Graph) (err error) {
		nodes := map[string]*cgraph.Node{}
		for name := range t.services {
			nodes[name], err = graph.CreateNode(name)
			if err != nil {
				return errors.Wrap(err, "failed to create graph node by graphviz")
			}
		}
		for _, svc := range t.services {
			if err := t.createEdges(graph, nodes, svc); err != nil {
				return err
			}
		}
		return nil
	})
}

func (t *Topology) ServiceGraphAsPNG(serviceName string) ([]byte, error) {
	return t.withGraph(graphviz.PNG, func(g *graphviz.Graphviz, graph *cgraph.Graph) error {
		nodes := map[string]*cgraph.Node{}
		if err := t.serviceGraph(g, graph, serviceName, nodes); err != nil {
			return err
		}
		return nil
	})
}

// TODO: refactoring
func (t *Topology) serviceGraph(g *graphviz.Graphviz, graph *cgraph.Graph, serviceName string, nodes map[string]*cgraph.Node) error {
	svc, ok := t.services[serviceName]
	if !ok {
		level.Warn(t.logger).Log("msg", "service not found", "service", serviceName)
		return nil
	}
	_, ok = nodes[serviceName]
	if ok {
		return nil
	}
	node, err := graph.CreateNode(serviceName)
	if err != nil {
		return errors.Wrap(err, "failed to create graph node by graphviz")
	}
	nodes[serviceName] = node

	for _, dep := range svc.Deps {
		if err := t.serviceGraph(g, graph, dep.Name, nodes); err != nil {
			return err
		}
	}
	if err := t.createEdges(graph, nodes, svc); err != nil {
		return err
	}
	return nil
}

func (t *Topology) createEdges(graph *cgraph.Graph, nodes map[string]*cgraph.Node, svc *Service) error {
	srcNode, ok := nodes[svc.Name]
	if !ok {
		level.Warn(t.logger).Log("msg", "service not found", "service", svc.Name)
		return nil
	}
	for _, dep := range svc.Deps {
		depName := dep.Name
		dstNode, ok := nodes[depName]
		if !ok {
			continue
		}
		_, err := graph.CreateEdge("", srcNode, dstNode)
		if err != nil {
			return errors.Wrap(err, "failed to create edge by graphviz")
		}
	}
	return nil
}

func (t *Topology) withGraph(format graphviz.Format, f GraphvizFunc) ([]byte, error) {
	g := graphviz.New()
	defer util.CloseWithLogOnErr(g)
	graph, err := g.Graph()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create graph by graphviz")
	}
	defer util.CloseWithLogOnErr(graph)
	graph.SetLayout(string(t.cfg.GetLayout()))

	buf := new(bytes.Buffer)
	if err := g.Render(graph, format, buf); err != nil {
		return nil, errors.Wrap(err, "failed to render by graphviz")
	}
	if err := f(g, graph); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
