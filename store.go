package modulyn

import (
	"sync"
)

type store struct {
	mu       sync.RWMutex
	features map[string]Feature
}

func (s *store) addOrUpdate(feature Feature) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.features == nil {
		s.features = make(map[string]Feature)
	}
	s.features[feature.Label] = feature
}

func (s *store) remove(feature Feature) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.features == nil {
		return
	}
	delete(s.features, feature.Label)
}
