package handler

import (
	"net/http"

	"logistics/internal/model"
	"logistics/pkg/httpx"
)

func (s *Server) registerParcelRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/parcels", s.createParcel)
	mux.HandleFunc("GET /api/parcels", s.listParcels)
	mux.HandleFunc("GET /api/parcels/{id}", s.getParcel)
	mux.HandleFunc("DELETE /api/parcels/{id}", s.deleteParcel)
	mux.HandleFunc("GET /api/track-events", s.listTrackEvents)
}

type parcelRequest struct {
	WaybillID   string  `json:"waybill_id"`
	Description string  `json:"description"`
	Weight      float64 `json:"weight"`
}

func (s *Server) createParcel(w http.ResponseWriter, r *http.Request) {
	var req parcelRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	p, err := s.svc.CreateParcel(model.Parcel{
		WaybillID: req.WaybillID, Description: req.Description, Weight: req.Weight,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, p)
}

func (s *Server) listParcels(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	items, total, err := s.svc.ListParcels(pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getParcel(w http.ResponseWriter, r *http.Request) {
	p, err := s.svc.GetParcel(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, p)
}

func (s *Server) deleteParcel(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteParcel(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) listTrackEvents(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.TrackEventFilter{
		WaybillID: r.URL.Query().Get("waybill_id"),
		Status:    r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListTrackEvents(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}
