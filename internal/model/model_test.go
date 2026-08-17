package model

import (
	"testing"
	"time"
)

func TestStationValidate(t *testing.T) {
	st := &Station{Name: "北京分拨", Address: "北京"}
	if err := st.Validate(); err != nil {
		t.Fatalf("valid: %v", err)
	}
	st.Name = ""
	if err := st.Validate(); err == nil {
		t.Fatalf("expect error for empty name")
	}
}

func TestWaybillValidate(t *testing.T) {
	w := &Waybill{Sender: "张三", Receiver: "李四", OriginStationID: "st1", DestStationID: "st2"}
	if err := w.Validate(); err != nil {
		t.Fatalf("valid: %v", err)
	}
	if w.Status != WaybillCreated {
		t.Fatalf("default status = %s", w.Status)
	}
	w.Sender = ""
	if err := w.Validate(); err == nil {
		t.Fatalf("expect error for empty sender")
	}
}

func TestWaybillTransitions(t *testing.T) {
	cases := []struct {
		from, to string
		want     bool
	}{
		{WaybillCreated, WaybillPicked, true},
		{WaybillCreated, WaybillDelivered, false},
		{WaybillPicked, WaybillInTransit, true},
		{WaybillInTransit, WaybillDelivering, true},
		{WaybillInTransit, WaybillException, true},
		{WaybillDelivering, WaybillDelivered, true},
		{WaybillDelivering, WaybillException, true},
		{WaybillException, WaybillInTransit, true},
		{WaybillException, WaybillDelivering, true},
		{WaybillDelivered, WaybillInTransit, false},
	}
	for _, c := range cases {
		if got := CanTransitionWaybill(c.from, c.to); got != c.want {
			t.Fatalf("transition %s->%s = %v, want %v", c.from, c.to, got, c.want)
		}
	}
}

func TestIsTerminalWaybill(t *testing.T) {
	if !IsTerminalWaybill(WaybillDelivered) {
		t.Fatalf("delivered should be terminal")
	}
	if IsTerminalWaybill(WaybillInTransit) {
		t.Fatalf("in_transit should not be terminal")
	}
}

func TestWaybillFilter(t *testing.T) {
	w := &Waybill{TrackingNo: "SF123", Receiver: "李四", Status: WaybillInTransit}
	if !(WaybillFilter{Status: WaybillInTransit}).Match(w) {
		t.Fatalf("expect match by status")
	}
	if !(WaybillFilter{Keyword: "SF123"}).Match(w) {
		t.Fatalf("expect match by tracking no")
	}
	if (WaybillFilter{Keyword: "不存在"}).Match(w) {
		t.Fatalf("expect no match")
	}
}

func TestParcelValidate(t *testing.T) {
	p := &Parcel{WaybillID: "w1", Weight: 2.5}
	if err := p.Validate(); err != nil {
		t.Fatalf("valid: %v", err)
	}
	p.Weight = -1
	if err := p.Validate(); err == nil {
		t.Fatalf("expect error for negative weight")
	}
}

func TestTrackEventValidate(t *testing.T) {
	te := &TrackEvent{WaybillID: "w1", CreatedAt: time.Now()}
	if err := te.Validate(); err != nil {
		t.Fatalf("valid: %v", err)
	}
	if te.Status != WaybillCreated {
		t.Fatalf("default status = %s", te.Status)
	}
	te.WaybillID = ""
	if err := te.Validate(); err == nil {
		t.Fatalf("expect error for empty waybill")
	}
}
