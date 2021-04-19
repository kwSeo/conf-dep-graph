package painter

import "github.com/goccy/go-graphviz"

type Config struct {
	Format string `yaml:"format"`
	Layout string `yaml:"layout"`
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
