package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"logistics/internal/config"
	"logistics/internal/model"
	"logistics/internal/service"
	"logistics/internal/store"
	"logistics/pkg/logger"
)

// newTestServer 装配一个带恢复中间件的完整 HTTP 服务，并预置 seed 条运单。
func newTestServer(t *testing.T, seed int) *Server {
	t.Helper()
	cfg := &config.Config{MaxPageSize: 100}
	log := logger.NewLevel(logger.LevelError)
	svc := service.New(store.NewMemoryStore(), log, cfg)
	from, err := svc.CreateStation(model.Station{Name: "北京", Address: "北京"})
	if err != nil {
		t.Fatalf("create station: %v", err)
	}
	to, err := svc.CreateStation(model.Station{Name: "上海", Address: "上海"})
	if err != nil {
		t.Fatalf("create station: %v", err)
	}
	for i := 0; i < seed; i++ {
		if _, err := svc.CreateWaybill(model.Waybill{
			Sender: "张三", Receiver: "李四",
			OriginStationID: from.ID, DestStationID: to.ID,
		}); err != nil {
			t.Fatalf("create waybill %d: %v", i, err)
		}
	}
	return NewServer(svc, log, cfg)
}

type apiResp struct {
	Code int             `json:"code"`
	Data json.RawMessage `json:"data"`
}

func listWaybills(t *testing.T, s *Server, query string) (int, []interface{}, int) {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/api/waybills?"+query, nil)
	rec := httptest.NewRecorder()
	s.Routes().ServeHTTP(rec, req)

	var resp apiResp
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v body=%s", err, rec.Body.String())
	}
	// data 形如 {"items":[...],"pagination":{"page":..,"size":..,"total":..}}
	var page struct {
		Items      []interface{} `json:"items"`
		Pagination struct {
			Total int `json:"total"`
		} `json:"pagination"`
	}
	_ = json.Unmarshal(resp.Data, &page)
	return rec.Code, page.Items, page.Pagination.Total
}

// 前端把 page 传成 0 或负数时，httpx 层已兜底为第 1 页，应返回第一页数据而非空。
func TestListWaybillsZeroOrNegativePage(t *testing.T) {
	s := newTestServer(t, 3)
	for _, q := range []string{"page=0", "page=-1", "page=-999&size=-5"} {
		code, items, total := listWaybills(t, s, q)
		if code != http.StatusOK {
			t.Fatalf("%s: expected 200, got %d", q, code)
		}
		if total != 3 || len(items) != 3 {
			t.Fatalf("%s: got len=%d total=%d, want 3 items total=3", q, len(items), total)
		}
	}
}

// 前端把 page 传成极大值（导致 (page-1)*size 整数溢出）时，绝不能 500/panic，
// 而应正常返回空页。
func TestListWaybillsHugePageNoPanic(t *testing.T) {
	s := newTestServer(t, 2)
	code, items, total := listWaybills(t, s, "page=9223372036854775807&size=20")
	if code != http.StatusOK {
		t.Fatalf("huge page should not panic: status=%d", code)
	}
	if total != 2 {
		t.Fatalf("total=%d want 2", total)
	}
	if len(items) != 0 {
		t.Fatalf("huge page should return empty page, got %d items", len(items))
	}
}
