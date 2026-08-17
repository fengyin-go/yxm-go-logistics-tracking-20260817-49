package service

import (
	"sort"
	"time"

	"logistics/internal/model"
	"logistics/internal/store"
	"logistics/pkg/idgen"
)

// CreateWaybill 创建运单：校验网点存在并生成运单号。
func (s *Service) CreateWaybill(input model.Waybill) (*model.Waybill, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetStation(input.OriginStationID); err != nil {
		return nil, err
	}
	if _, err := s.store.GetStation(input.DestStationID); err != nil {
		return nil, err
	}
	now := time.Now()
	w := &model.Waybill{
		ID:               idgen.Hex(),
		TrackingNo:       "SF" + idgen.HexN(6),
		Sender:           input.Sender,
		Receiver:         input.Receiver,
		OriginStationID:  input.OriginStationID,
		DestStationID:    input.DestStationID,
		CurrentStationID: input.OriginStationID,
		Status:           model.WaybillCreated,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	if err := s.store.CreateWaybill(w); err != nil {
		return nil, err
	}
	// 记录初始轨迹
	s.appendTrack(w.ID, w.CurrentStationID, w.Status, "运单已创建")
	return w, nil
}

func (s *Service) GetWaybill(id string) (*model.Waybill, error) {
	return s.store.GetWaybill(id)
}

func (s *Service) GetWaybillByTrackingNo(no string) (*model.Waybill, error) {
	return s.store.GetWaybillByTrackingNo(no)
}

func (s *Service) ListWaybills(filter model.WaybillFilter, page, size int) ([]*model.Waybill, int, error) {
	all := s.store.ListWaybills()
	matched := make([]*model.Waybill, 0, len(all))
	for _, w := range all {
		if filter.Match(w) {
			matched = append(matched, w)
		}
	}
	sort.Slice(matched, func(i, j int) bool { return matched[i].CreatedAt.After(matched[j].CreatedAt) })
	total := len(matched)
	start, end := pageBounds(total, page, size)
	return matched[start:end], total, nil
}

// Transition 流转运单状态并追加轨迹。
func (s *Service) Transition(id, targetStatus, stationID, description string) (*model.Waybill, error) {
	w, err := s.store.GetWaybill(id)
	if err != nil {
		return nil, err
	}
	if !model.CanTransitionWaybill(w.Status, targetStatus) {
		return nil, store.ErrConflict
	}
	if stationID != "" {
		if _, err := s.store.GetStation(stationID); err != nil {
			return nil, err
		}
		w.CurrentStationID = stationID
	}
	now := time.Now()
	w.Status = targetStatus
	w.UpdatedAt = now
	if targetStatus == model.WaybillDelivered {
		w.DeliveredAt = &now
	}
	if err := s.store.UpdateWaybill(w); err != nil {
		return nil, err
	}
	s.appendTrack(w.ID, w.CurrentStationID, targetStatus, description)
	return w, nil
}

// appendTrack 追加一条轨迹记录。
func (s *Service) appendTrack(waybillID, stationID, status, description string) {
	t := &model.TrackEvent{
		ID:          idgen.Hex(),
		WaybillID:   waybillID,
		StationID:   stationID,
		Status:      status,
		Description: description,
		CreatedAt:   time.Now(),
	}
	_ = s.store.CreateTrackEvent(t)
}

// Track 查询运单轨迹（按时间升序）。
func (s *Service) Track(waybillID string) ([]*model.TrackEvent, error) {
	if _, err := s.store.GetWaybill(waybillID); err != nil {
		return nil, err
	}
	events := make([]*model.TrackEvent, 0)
	for _, t := range s.store.ListTrackEvents() {
		if t.WaybillID == waybillID {
			events = append(events, t)
		}
	}
	sort.Slice(events, func(i, j int) bool { return events[i].CreatedAt.Before(events[j].CreatedAt) })
	return events, nil
}
