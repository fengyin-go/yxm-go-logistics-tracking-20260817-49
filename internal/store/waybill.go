package store

import (
	"logistics/internal/model"
)

func (s *MemoryStore) CreateWaybill(w *model.Waybill) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.waybills {
		if exist.TrackingNo == w.TrackingNo {
			return ErrConflict
		}
	}
	s.waybills[w.ID] = w
	return nil
}

func (s *MemoryStore) GetWaybill(id string) (*model.Waybill, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	w, ok := s.waybills[id]
	if !ok {
		return nil, ErrNotFound
	}
	return w, nil
}

func (s *MemoryStore) GetWaybillByTrackingNo(no string) (*model.Waybill, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, w := range s.waybills {
		if w.TrackingNo == no {
			return w, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListWaybills() []*model.Waybill {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Waybill, 0, len(s.waybills))
	for _, w := range s.waybills {
		list = append(list, w)
	}
	return list
}

func (s *MemoryStore) UpdateWaybill(w *model.Waybill) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.waybills[w.ID]; !ok {
		return ErrNotFound
	}
	s.waybills[w.ID] = w
	return nil
}

func (s *MemoryStore) DeleteWaybill(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.waybills[id]; !ok {
		return ErrNotFound
	}
	delete(s.waybills, id)
	return nil
}
