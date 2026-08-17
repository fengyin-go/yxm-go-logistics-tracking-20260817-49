package service

import (
	"sort"
	"time"

	"logistics/internal/model"
	"logistics/pkg/idgen"
)

func (s *Service) CreateStation(input model.Station) (*model.Station, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	st := &model.Station{
		ID:        idgen.Hex(),
		Name:      input.Name,
		Address:   input.Address,
		Phone:     input.Phone,
		CreatedAt: time.Now(),
	}
	if err := s.store.CreateStation(st); err != nil {
		return nil, err
	}
	return st, nil
}

func (s *Service) GetStation(id string) (*model.Station, error) {
	return s.store.GetStation(id)
}

func (s *Service) ListStations(page, size int) ([]*model.Station, int, error) {
	all := s.store.ListStations()
	sort.Slice(all, func(i, j int) bool { return all[i].CreatedAt.After(all[j].CreatedAt) })
	total := len(all)
	start := (page - 1) * size
	if start >= total {
		return []*model.Station{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return all[start:end], total, nil
}

func (s *Service) UpdateStation(id string, input model.Station) (*model.Station, error) {
	existing, err := s.store.GetStation(id)
	if err != nil {
		return nil, err
	}
	if input.Name != "" {
		existing.Name = input.Name
	}
	if input.Address != "" {
		existing.Address = input.Address
	}
	if input.Phone != "" {
		existing.Phone = input.Phone
	}
	if err := existing.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateStation(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *Service) DeleteStation(id string) error {
	// 跨实体校验：网点被运单引用时不允许删除
	for _, w := range s.store.ListWaybills() {
		if w.OriginStationID == id || w.DestStationID == id || w.CurrentStationID == id {
			return model.NewValidationError("station", "该网点已被运单引用，无法删除")
		}
	}
	return s.store.DeleteStation(id)
}
