// Package store 定义数据访问接口与内存实现。
package store

import (
	"errors"

	"logistics/internal/model"
)

var (
	ErrNotFound = errors.New("记录不存在")
	ErrConflict = errors.New("记录已存在或状态冲突")
)

// Store 聚合全部实体的数据访问方法。
type Store interface {
	// 网点
	CreateStation(st *model.Station) error
	GetStation(id string) (*model.Station, error)
	ListStations() []*model.Station
	UpdateStation(st *model.Station) error
	DeleteStation(id string) error

	// 运单
	CreateWaybill(w *model.Waybill) error
	GetWaybill(id string) (*model.Waybill, error)
	GetWaybillByTrackingNo(no string) (*model.Waybill, error)
	ListWaybills() []*model.Waybill
	UpdateWaybill(w *model.Waybill) error
	DeleteWaybill(id string) error

	// 包裹
	CreateParcel(p *model.Parcel) error
	GetParcel(id string) (*model.Parcel, error)
	GetParcelByWaybill(waybillID string) (*model.Parcel, error)
	ListParcels() []*model.Parcel
	DeleteParcel(id string) error

	// 轨迹
	CreateTrackEvent(t *model.TrackEvent) error
	GetTrackEvent(id string) (*model.TrackEvent, error)
	ListTrackEvents() []*model.TrackEvent
}
