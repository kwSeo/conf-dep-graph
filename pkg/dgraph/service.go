package dgraph

import (
	"strings"
	"sync"
)

type Service struct {
	Name     string
	Deps     []*Service
	Contents []ContentFile

	depsLock sync.Mutex
}

func NewService(name string) *Service {
	return &Service{
		Name:     name,
		depsLock: sync.Mutex{},
	}
}

func (s *Service) DependOn(other *Service) bool {
	for _, content := range s.Contents {
		if strings.Contains(content.String(), other.Name) {
			return true
		}
	}
	return false
}

func (s *Service) AddDep(other *Service) {
	s.depsLock.Lock()
	defer s.depsLock.Unlock()
	s.Deps = append(s.Deps, other)
}

type ContentFile struct {
	URL     string
	content []byte
}

func (c *ContentFile) String() string {
	return string(c.content)
}

func (c *ContentFile) Bytes() []byte {
	return c.content
}
