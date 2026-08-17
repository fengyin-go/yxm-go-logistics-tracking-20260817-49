package store

import (
	"logistics/internal/model"
)

func (s *MemoryStore) CreateParcel(p *model.Parcel) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.parcels {
		if exist.WaybillID == p.WaybillID {
			return ErrConflict
		}
	}
	s.parcels[p.ID] = p
	return nil
}

func (s *MemoryStore) GetParcel(id string) (*model.Parcel, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.parcels[id]
	if !ok {
		return nil, ErrNotFound
	}
	return p, nil
}

func (s *MemoryStore) GetParcelByWaybill(waybillID string) (*model.Parcel, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, p := range s.parcels {
		if p.WaybillID == waybillID {
			return p, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListParcels() []*model.Parcel {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Parcel, 0, len(s.parcels))
	for _, p := range s.parcels {
		list = append(list, p)
	}
	return list
}

func (s *MemoryStore) DeleteParcel(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.parcels[id]; !ok {
		return ErrNotFound
	}
	delete(s.parcels, id)
	return nil
}
