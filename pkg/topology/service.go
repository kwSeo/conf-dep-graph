package topology

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/kwseo/conf-dep-graph/pkg/util"
)

type Service struct {
	Name       string     `yaml:"name"`
	Keyword    string     `yaml:"keyword"`
	TargetFile string     `yaml:"target_file"`
	Deps       []*Service `yaml:"deps"`

	Content []byte `yaml:"-"`
}

func (s *Service) LoadTargetFile() (content []byte, err error) {
	if s.TargetFile == "" {
		content = []byte{}
	} else {
		switch pathTypeOf(s.TargetFile) {
		case FILESYSTEM:
			content, err = readFile(s.TargetFile)
		case URL:
			content, err = readURL(s.TargetFile)
		}
	}
	s.Content = content
	return
}

func (s *Service) DependOn(other *Service) bool {
	return strings.Contains(string(s.Content), other.Keyword)
}

func (s *Service) AddDependency(other *Service) {
	s.Deps = append(s.Deps, other)
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
