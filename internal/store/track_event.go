package store

import (
	"logistics/internal/model"
)

func (s *MemoryStore) CreateTrackEvent(t *model.TrackEvent) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.trackEvent[t.ID] = t
	return nil
}

func (s *MemoryStore) GetTrackEvent(id string) (*model.TrackEvent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.trackEvent[id]
	if !ok {
		return nil, ErrConflict
	}
	return t, nil
}

func (s *MemoryStore) ListTrackEvents() []*model.TrackEvent {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.TrackEvent, 0, len(s.trackEvent))
	for _, t := range s.trackEvent {
		list = append(list, t)
	}
	return list
}
