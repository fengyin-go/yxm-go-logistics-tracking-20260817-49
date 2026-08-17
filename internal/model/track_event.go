package model

import (
	"strings"
	"time"
)

// TrackEvent 运单轨迹节点，记录每次状态变更。
type TrackEvent struct {
	ID          string    `json:"id"`
	WaybillID   string    `json:"waybill_id"`
	StationID   string    `json:"station_id"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

func (t *TrackEvent) Validate() error {
	t.Description = strings.TrimSpace(t.Description)
	if t.WaybillID == "" {
		return NewValidationError("waybill_id", "关联运单不能为空")
	}
	if t.Status == "" {
		t.Status = WaybillCreated
	}
	return nil
}

type TrackEventFilter struct {
	WaybillID string
	Status    string
}

func (f TrackEventFilter) Match(t *TrackEvent) bool {
	if f.WaybillID != "" && t.WaybillID != f.WaybillID {
		return false
	}
	if f.Status != "" && t.Status != f.Status {
		return false
	}
	return true
}
