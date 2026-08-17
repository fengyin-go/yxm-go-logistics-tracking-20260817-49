package model

import (
	"strings"
	"time"
)

// Parcel 包裹，一个运单可挂一个包裹信息。
type Parcel struct {
	ID          string    `json:"id"`
	WaybillID   string    `json:"waybill_id"`
	Description string    `json:"description"`
	Weight      float64   `json:"weight"`
	CreatedAt   time.Time `json:"created_at"`
}

func (p *Parcel) Validate() error {
	p.Description = strings.TrimSpace(p.Description)
	if p.WaybillID == "" {
		return NewValidationError("waybill_id", "关联运单不能为空")
	}
	if p.Weight < 0 {
		return NewValidationError("weight", "重量不能为负")
	}
	return nil
}
