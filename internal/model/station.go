package model

import (
	"strings"
	"time"
)

// Station 物流网点。
type Station struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
}

func (st *Station) Validate() error {
	st.Name = strings.TrimSpace(st.Name)
	st.Address = strings.TrimSpace(st.Address)
	if st.Name == "" {
		return NewValidationError("name", "网点名称不能为空")
	}
	if st.Address == "" {
		return NewValidationError("address", "网点地址不能为空")
	}
	return nil
}
