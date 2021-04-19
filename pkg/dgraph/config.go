package dgraph

import (
	"io/ioutil"

	"github.com/kwseo/conf-dep-graph/pkg/painter"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Config struct {
	OutputConfig   painter.Config  `yaml:"output"`
	ServiceConfigs []ServiceConfig `yaml:"services"`
}

type ServiceConfig struct {
	Name         string   `yaml:"name"`
	Keywords     []string `yaml:"keywords"`
	ContentFiles []string `yaml:"content_files"`
	Deps         []string `yaml:"deps"`
}

func LoadConfig(filePath string) (*Config, error) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read file: %s", filePath)
	}
	var cfg Config
	if err := yaml.Unmarshal(file, cfg); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal YAML config file")
	}
	return &cfg, nil
}
