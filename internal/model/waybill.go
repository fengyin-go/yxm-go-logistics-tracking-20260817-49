package model

import (
	"strings"
	"time"
)

const (
	WaybillCreated    = "created"
	WaybillPicked     = "picked"
	WaybillInTransit  = "in_transit"
	WaybillDelivering = "delivering"
	WaybillDelivered  = "delivered"
	WaybillException  = "exception"
)

// waybillTransitions 运单状态机。
var waybillTransitions = map[string]map[string]bool{
	WaybillCreated:    {WaybillPicked: true},
	WaybillPicked:     {WaybillInTransit: true},
	WaybillInTransit:  {WaybillDelivering: true, WaybillException: true},
	WaybillDelivering: {WaybillDelivered: true, WaybillException: true},
	WaybillException:  {WaybillInTransit: true},
}

// CanTransitionWaybill 判断运单状态流转是否合法。
func CanTransitionWaybill(from, to string) bool {
	if m, ok := waybillTransitions[from]; ok {
		return m[to]
	}
	return false
}

// IsTerminalWaybill 判断是否为终态。
func IsTerminalWaybill(status string) bool {
	return status == WaybillDelivered
}

// Waybill 运单。
type Waybill struct {
	ID               string     `json:"id"`
	TrackingNo       string     `json:"tracking_no"`
	Sender           string     `json:"sender"`
	Receiver         string     `json:"receiver"`
	OriginStationID  string     `json:"origin_station_id"`
	DestStationID    string     `json:"dest_station_id"`
	CurrentStationID string     `json:"current_station_id"`
	Status           string     `json:"status"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeliveredAt      *time.Time `json:"delivered_at,omitempty"`
}

func (w *Waybill) Validate() error {
	w.Sender = strings.TrimSpace(w.Sender)
	w.Receiver = strings.TrimSpace(w.Receiver)
	if w.Sender == "" {
		return NewValidationError("sender", "寄件人不能为空")
	}
	if w.Receiver == "" {
		return NewValidationError("receiver", "收件人不能为空")
	}
	if w.OriginStationID == "" {
		return NewValidationError("origin_station_id", "始发网点不能为空")
	}
	if w.DestStationID == "" {
		return NewValidationError("dest_station_id", "目的网点不能为空")
	}
	if w.Status == "" {
		w.Status = WaybillCreated
	}
	if w.Status != WaybillCreated && w.Status != WaybillPicked && w.Status != WaybillInTransit &&
		w.Status != WaybillDelivering && w.Status != WaybillDelivered && w.Status != WaybillException {
		return NewValidationError("status", "运单状态不合法")
	}
	return nil
}

type WaybillFilter struct {
	Status  string
	Keyword string
}

func (f WaybillFilter) Match(w *Waybill) bool {
	if f.Status != "" && w.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(w.TrackingNo), k) &&
			!strings.Contains(strings.ToLower(w.Receiver), k) {
			return false
		}
	}
	return true
}
