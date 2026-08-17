package store

import (
	"testing"
	"time"

	"logistics/internal/model"
)

func newTestStore() *MemoryStore { return NewMemoryStore() }

func TestStationCRUD(t *testing.T) {
	s := newTestStore()
	st := &model.Station{ID: "st1", Name: "北京分拨中心", Address: "北京市大兴区", CreatedAt: time.Now()}
	if err := s.CreateStation(st); err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := s.CreateStation(&model.Station{ID: "st2", Name: "北京分拨中心", Address: "x"}); err != ErrConflict {
		t.Fatalf("expect conflict, got %v", err)
	}
	got, err := s.GetStation("st1")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Name != "北京分拨中心" {
		t.Fatalf("name = %s", got.Name)
	}
	if _, err := s.GetStation("missing"); err != ErrNotFound {
		t.Fatalf("expect not found, got %v", err)
	}
	if n := len(s.ListStations()); n != 1 {
		t.Fatalf("list len = %d", n)
	}
	got.Phone = "010-12345678"
	if err := s.UpdateStation(got); err != nil {
		t.Fatalf("update: %v", err)
	}
	if err := s.DeleteStation("st1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
}

func TestWaybillCRUD(t *testing.T) {
	s := newTestStore()
	w := &model.Waybill{ID: "w1", TrackingNo: "SF123", Sender: "张三", Receiver: "李四", OriginStationID: "st1", DestStationID: "st2", Status: model.WaybillCreated}
	if err := s.CreateWaybill(w); err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := s.CreateWaybill(&model.Waybill{ID: "w2", TrackingNo: "SF123", Sender: "a", Receiver: "b"}); err != ErrConflict {
		t.Fatalf("expect conflict by tracking no, got %v", err)
	}
	got, _ := s.GetWaybill("w1")
	got.Status = model.WaybillPicked
	if err := s.UpdateWaybill(got); err != nil {
		t.Fatalf("update: %v", err)
	}
	byNo, _ := s.GetWaybillByTrackingNo("SF123")
	if byNo.Status != model.WaybillPicked {
		t.Fatalf("by no status = %s", byNo.Status)
	}
	if err := s.DeleteWaybill("w1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
}

func TestParcelCRUD(t *testing.T) {
	s := newTestStore()
	p := &model.Parcel{ID: "p1", WaybillID: "w1", Description: "电子产品", Weight: 1.5}
	if err := s.CreateParcel(p); err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := s.CreateParcel(&model.Parcel{ID: "p2", WaybillID: "w1"}); err != ErrConflict {
		t.Fatalf("expect conflict by waybill, got %v", err)
	}
	if _, err := s.GetParcel("p1"); err != nil {
		t.Fatalf("get: %v", err)
	}
	if _, err := s.GetParcelByWaybill("w1"); err != nil {
		t.Fatalf("get by waybill: %v", err)
	}
	if err := s.DeleteParcel("p1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
}

func TestTrackEventCRUD(t *testing.T) {
	s := newTestStore()
	t1 := &model.TrackEvent{ID: "t1", WaybillID: "w1", Status: model.WaybillCreated, CreatedAt: time.Now()}
	if err := s.CreateTrackEvent(t1); err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, err := s.GetTrackEvent("t1"); err != nil {
		t.Fatalf("get: %v", err)
	}
	if n := len(s.ListTrackEvents()); n != 1 {
		t.Fatalf("list len = %d", n)
	}
}

func TestParcelListDoesNotExposeStoredPointers(t *testing.T) {
	s := newTestStore()
	p := &model.Parcel{ID: "p1", WaybillID: "w1", Description: "电子产品", Weight: 1.5}
	if err := s.CreateParcel(p); err != nil {
		t.Fatalf("create: %v", err)
	}
	items := s.ListParcels()
	if len(items) != 1 {
		t.Fatalf("list len = %d", len(items))
	}
	items[0].Description = "被外部污染"
	got, err := s.GetParcel("p1")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Description != "电子产品" {
		t.Fatalf("stored parcel was mutated through list result: %q", got.Description)
	}
}
