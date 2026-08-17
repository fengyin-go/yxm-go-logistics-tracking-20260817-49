package service

import (
	"testing"

	"logistics/internal/config"
	"logistics/internal/model"
	"logistics/internal/store"
	"logistics/pkg/logger"
)

func newTestService() *Service {
	cfg := &config.Config{MaxPageSize: 100}
	log := logger.NewLevel(logger.LevelError)
	return New(store.NewMemoryStore(), log, cfg)
}

func setupStations(t *testing.T, s *Service) (string, string) {
	t.Helper()
	from, err := s.CreateStation(model.Station{Name: "北京分拨", Address: "北京"})
	if err != nil {
		t.Fatalf("create from: %v", err)
	}
	to, err := s.CreateStation(model.Station{Name: "上海分拨", Address: "上海"})
	if err != nil {
		t.Fatalf("create to: %v", err)
	}
	return from.ID, to.ID
}

func TestWaybillFullFlow(t *testing.T) {
	s := newTestService()
	from, to := setupStations(t, s)

	wb, err := s.CreateWaybill(model.Waybill{Sender: "张三", Receiver: "李四", OriginStationID: from, DestStationID: to})
	if err != nil {
		t.Fatalf("create waybill: %v", err)
	}
	if wb.Status != model.WaybillCreated || wb.TrackingNo == "" {
		t.Fatalf("status=%s no=%s", wb.Status, wb.TrackingNo)
	}

	// 揽收
	wb, err = s.Transition(wb.ID, model.WaybillPicked, from, "已揽收")
	if err != nil || wb.Status != model.WaybillPicked {
		t.Fatalf("pick: %v status=%s", err, wb.Status)
	}
	// 运输中
	wb, _ = s.Transition(wb.ID, model.WaybillInTransit, "", "发往上海")
	// 派送中
	wb, _ = s.Transition(wb.ID, model.WaybillDelivering, to, "到达派送")
	// 签收
	wb, err = s.Transition(wb.ID, model.WaybillDelivered, to, "已签收")
	if err != nil || wb.Status != model.WaybillDelivered || wb.DeliveredAt == nil {
		t.Fatalf("deliver: %v status=%s", err, wb.Status)
	}

	// 轨迹应有 5 条（创建 + 4 次流转）
	events, err := s.Track(wb.ID)
	if err != nil || len(events) != 5 {
		t.Fatalf("track: %v len=%d", err, len(events))
	}
}

func TestWaybillInvalidTransition(t *testing.T) {
	s := newTestService()
	from, to := setupStations(t, s)
	wb, _ := s.CreateWaybill(model.Waybill{Sender: "张三", Receiver: "李四", OriginStationID: from, DestStationID: to})
	// created 不能直接到 delivered
	if _, err := s.Transition(wb.ID, model.WaybillDelivered, to, ""); err == nil {
		t.Fatalf("expect error for invalid transition")
	}
}

func TestCreateWaybillInvalidStation(t *testing.T) {
	s := newTestService()
	if _, err := s.CreateWaybill(model.Waybill{Sender: "a", Receiver: "b", OriginStationID: "missing", DestStationID: "x"}); err == nil {
		t.Fatalf("expect error for missing station")
	}
}

func TestDeleteStationBlockedByWaybill(t *testing.T) {
	s := newTestService()
	from, to := setupStations(t, s)
	s.CreateWaybill(model.Waybill{Sender: "张三", Receiver: "李四", OriginStationID: from, DestStationID: to})
	if err := s.DeleteStation(from); err == nil {
		t.Fatalf("expect error deleting referenced station")
	}
}

func TestParcelService(t *testing.T) {
	s := newTestService()
	from, to := setupStations(t, s)
	wb, _ := s.CreateWaybill(model.Waybill{Sender: "张三", Receiver: "李四", OriginStationID: from, DestStationID: to})
	p, err := s.CreateParcel(model.Parcel{WaybillID: wb.ID, Description: "电子产品", Weight: 2.5})
	if err != nil {
		t.Fatalf("create parcel: %v", err)
	}
	if _, err := s.CreateParcel(model.Parcel{WaybillID: wb.ID}); err == nil {
		t.Fatalf("expect conflict on duplicate parcel")
	}
	if _, err := s.GetParcel(p.ID); err != nil {
		t.Fatalf("get parcel: %v", err)
	}
}

func TestExceptionWaybillCanResumeDeliveryWithCompleteTrack(t *testing.T) {
	s := newTestService()
	from, to := setupStations(t, s)
	wb, err := s.CreateWaybill(model.Waybill{Sender: "张三", Receiver: "李四", OriginStationID: from, DestStationID: to})
	if err != nil {
		t.Fatalf("create waybill: %v", err)
	}
	if _, err = s.Transition(wb.ID, model.WaybillPicked, from, "已揽收"); err != nil {
		t.Fatalf("picked: %v", err)
	}
	if _, err = s.Transition(wb.ID, model.WaybillInTransit, from, "干线运输"); err != nil {
		t.Fatalf("in transit: %v", err)
	}
	if _, err = s.Transition(wb.ID, model.WaybillException, from, "分拣异常"); err != nil {
		t.Fatalf("exception: %v", err)
	}
	if _, err = s.Transition(wb.ID, model.WaybillDelivering, to, "异常解除后派送"); err != nil {
		t.Fatalf("resume delivering after exception: %v", err)
	}
	events, err := s.Track(wb.ID)
	if err != nil {
		t.Fatalf("track: %v", err)
	}
	if len(events) != 5 {
		t.Fatalf("track len = %d, want 5", len(events))
	}
	if events[3].Status != model.WaybillException || events[3].Description != "分拣异常" {
		t.Fatalf("exception event lost: %#v", events[3])
	}
	if events[4].Status != model.WaybillDelivering || events[4].StationID != to {
		t.Fatalf("resume event = %#v", events[4])
	}
}
