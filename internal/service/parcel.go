package service

import (
	"sort"
	"time"

	"logistics/internal/model"
	"logistics/pkg/idgen"
)

func (s *Service) CreateParcel(input model.Parcel) (*model.Parcel, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetWaybill(input.WaybillID); err != nil {
		return nil, err
	}
	p := &model.Parcel{
		ID:          idgen.Hex(),
		WaybillID:   input.WaybillID,
		Description: input.Description,
		Weight:      input.Weight,
		CreatedAt:   time.Now(),
	}
	if err := s.store.CreateParcel(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) GetParcel(id string) (*model.Parcel, error) {
	return s.store.GetParcel(id)
}

func (s *Service) ListParcels(page, size int) ([]*model.Parcel, int, error) {
	all := s.store.ListParcels()
	sort.Slice(all, func(i, j int) bool { return all[i].CreatedAt.After(all[j].CreatedAt) })
	total := len(all)
	start, end := pageBounds(total, page, size)
	return all[start:end], total, nil
}

func (s *Service) DeleteParcel(id string) error {
	return s.store.DeleteParcel(id)
}

func (s *Service) ListTrackEvents(filter model.TrackEventFilter, page, size int) ([]*model.TrackEvent, int, error) {
	all := s.store.ListTrackEvents()
	matched := make([]*model.TrackEvent, 0, len(all))
	for _, t := range all {
		if filter.Match(t) {
			matched = append(matched, t)
		}
	}
	sort.Slice(matched, func(i, j int) bool { return matched[i].CreatedAt.After(matched[j].CreatedAt) })
	total := len(matched)
	start, end := pageBounds(total, page, size)
	return matched[start:end], total, nil
}
