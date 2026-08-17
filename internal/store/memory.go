package store

import (
	"sync"

	"logistics/internal/model"
)

// MemoryStore 基于内存的 Store 实现，使用读写锁保证并发安全。
type MemoryStore struct {
	mu         sync.RWMutex
	stations   map[string]*model.Station
	waybills   map[string]*model.Waybill
	parcels    map[string]*model.Parcel
	trackEvent map[string]*model.TrackEvent
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		stations:   make(map[string]*model.Station),
		waybills:   make(map[string]*model.Waybill),
		parcels:    make(map[string]*model.Parcel),
		trackEvent: make(map[string]*model.TrackEvent),
	}
}

var _ Store = (*MemoryStore)(nil)
