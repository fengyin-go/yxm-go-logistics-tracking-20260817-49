package handler

import (
	"net/http"

	"logistics/internal/model"
	"logistics/pkg/httpx"
)

func (s *Server) registerStationRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/stations", s.createStation)
	mux.HandleFunc("GET /api/stations", s.listStations)
	mux.HandleFunc("GET /api/stations/{id}", s.getStation)
	mux.HandleFunc("PUT /api/stations/{id}", s.updateStation)
	mux.HandleFunc("DELETE /api/stations/{id}", s.deleteStation)
}

type stationRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
}

func (s *Server) createStation(w http.ResponseWriter, r *http.Request) {
	var req stationRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	st, err := s.svc.CreateStation(model.Station{Name: req.Name, Address: req.Address, Phone: req.Phone})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, st)
}

func (s *Server) listStations(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	items, total, err := s.svc.ListStations(pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getStation(w http.ResponseWriter, r *http.Request) {
	st, err := s.svc.GetStation(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, st)
}

func (s *Server) updateStation(w http.ResponseWriter, r *http.Request) {
	var req stationRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	st, err := s.svc.UpdateStation(r.PathValue("id"), model.Station{Name: req.Name, Address: req.Address, Phone: req.Phone})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, st)
}

func (s *Server) deleteStation(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteStation(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
