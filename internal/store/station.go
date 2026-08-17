package store

import (
	"logistics/internal/model"
)

func (s *MemoryStore) CreateStation(st *model.Station) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.stations {
		if exist.Name == st.Name {
			return ErrConflict
		}
	}
	s.stations[st.ID] = st
	return nil
}

func (s *MemoryStore) GetStation(id string) (*model.Station, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	st, ok := s.stations[id]
	if !ok {
		return nil, ErrNotFound
	}
	return st, nil
}

func (s *MemoryStore) ListStations() []*model.Station {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Station, 0, len(s.stations))
	for _, st := range s.stations {
		list = append(list, st)
	}
	return list
}

func (s *MemoryStore) UpdateStation(st *model.Station) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.stations[st.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.stations {
		if exist.ID != st.ID && exist.Name == st.Name {
			return ErrConflict
		}
	}
	s.stations[st.ID] = st
	return nil
}

func (s *MemoryStore) DeleteStation(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.stations[id]; !ok {
		return ErrNotFound
	}
	delete(s.stations, id)
	return nil
}
