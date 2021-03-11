package topology

import (
	"bytes"
	"io/ioutil"
	"log"

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
}

func New(cfg *Config) *Topology {
	return &Topology{
		cfg:      cfg,
		services: map[string]*Service{},
	}
}

func (t *Topology) AddService(newSvc *Service) {
	if _, exist := t.services[newSvc.Name]; exist {
		log.Println("Already existed service")
		return
	}
	for _, svc := range t.services {
		if svc.DependOn(newSvc) {
			svc.AddDependency(*newSvc)
		}
		if newSvc.DependOn(svc) {
			newSvc.AddDependency(*svc)
		}
	}
	t.services[newSvc.Name] = newSvc
}

func (t *Topology) GraphAsPNG() ([]byte, error) {
	g := graphviz.New()
	defer util.CloseWithLogOnErr(g)
	graph, err := g.Graph()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create graph by graphviz")
	}
	defer util.CloseWithLogOnErr(graph)
	graph.SetLayout(string(t.cfg.GetLayout()))

	nodes := map[string]*cgraph.Node{}
	for name := range t.services {
		nodes[name], err = graph.CreateNode(name)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create graph node by graphviz")
		}
	}
	for _, svc := range t.services {
		srcNode, ok := nodes[svc.Name]
		if !ok {
			continue
		}
		for _, depName := range svc.Deps {
			dstNode, ok := nodes[depName]
			if !ok {
				continue
			}
			_, err := graph.CreateEdge("", srcNode, dstNode)
			if err != nil {
				return nil, errors.Wrap(err, "failed to create edge by graphviz")
			}
		}
	}
	buf := new(bytes.Buffer)
	if err := g.Render(graph, graphviz.PNG, buf); err != nil {
		return nil, errors.Wrap(err, "failed to render by graphviz")
	}
	return buf.Bytes(), nil
}

// func (t *Topology) CreateServices() ([]*Service, error) {
// 	services := map[string]*Service{}
// 	for _, svcCfg := range t.cfg.Services {
// 		content, err := read(svcCfg.TargetFile)
// 		if err != nil {
// 			return nil, errors.Wrap(err, "failed to read target file")
// 		}
// 	}
// 	return services, nil
// }
