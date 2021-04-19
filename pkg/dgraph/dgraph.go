package dgraph

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/kwseo/conf-dep-graph/pkg/util"
)

// DGraph indicate a dependency graph.
type DGraph struct {
	services map[string]*Service
}

func New(cfg *Config) *DGraph {
	dg := &DGraph{
		services: map[string]*Service{},
	}
	initDGraph(dg, cfg)
	return dg
}

func initDGraph(dg *DGraph, cfg *Config) error {
	for _, serviceConfig := range cfg.ServiceConfigs {
		service := NewService(serviceConfig.Name)
		// init content files
		for _, contentFilePath := range serviceConfig.ContentFiles {
			content, err := readContent(contentFilePath)
			if err != nil {
				return err
			}
			contentFile := ContentFile{
				URL:     contentFilePath,
				content: content,
			}
			service.Contents = append(service.Contents, contentFile)
		}
		dg.services[service.Name] = service
	}
	// init predefined depdendencies
	panic("implement me")
}

type PathType int

const (
	FILESYSTEM PathType = iota
	URL
)

var ErrUnexpectedContentPath = errors.New("unexpected content file path")

func readContent(path string) ([]byte, error) {
	switch pathTypeOf(path) {
	case URL:
		return readURL(path)
	case FILESYSTEM:
		return readFile(path)
	default:
		return nil, ErrUnexpectedContentPath
	}
}

func pathTypeOf(path string) PathType {
	if strings.HasPrefix(path, "http") {
		return URL
	} else {
		return FILESYSTEM
	}
}

func readFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func readURL(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer util.CloseWithLogOnErr(resp.Body)
	return ioutil.ReadAll(resp.Body)
}
