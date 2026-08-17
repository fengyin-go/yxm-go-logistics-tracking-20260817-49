package handler

import (
	"net/http"

	"logistics/internal/model"
	"logistics/pkg/httpx"
)

func (s *Server) registerWaybillRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/waybills", s.createWaybill)
	mux.HandleFunc("GET /api/waybills", s.listWaybills)
	mux.HandleFunc("GET /api/waybills/{id}", s.getWaybill)
	mux.HandleFunc("POST /api/waybills/{id}/transition", s.transitionWaybill)
	mux.HandleFunc("GET /api/waybills/{id}/track", s.trackWaybill)
}

type waybillRequest struct {
	Sender          string `json:"sender"`
	Receiver        string `json:"receiver"`
	OriginStationID string `json:"origin_station_id"`
	DestStationID   string `json:"dest_station_id"`
}

func (s *Server) createWaybill(w http.ResponseWriter, r *http.Request) {
	var req waybillRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	wb, err := s.svc.CreateWaybill(model.Waybill{
		Sender: req.Sender, Receiver: req.Receiver,
		OriginStationID: req.OriginStationID, DestStationID: req.DestStationID,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, wb)
}

func (s *Server) listWaybills(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.WaybillFilter{
		Status:  r.URL.Query().Get("status"),
		Keyword: r.URL.Query().Get("q"),
	}
	items, total, err := s.svc.ListWaybills(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getWaybill(w http.ResponseWriter, r *http.Request) {
	wb, err := s.svc.GetWaybill(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, wb)
}

type transitionRequest struct {
	Status      string `json:"status"`
	StationID   string `json:"station_id"`
	Description string `json:"description"`
}

func (s *Server) transitionWaybill(w http.ResponseWriter, r *http.Request) {
	var req transitionRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	wb, err := s.svc.Transition(r.PathValue("id"), req.Status, req.StationID, req.Description)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, wb)
}

func (s *Server) trackWaybill(w http.ResponseWriter, r *http.Request) {
	events, err := s.svc.Track(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, events)
}
