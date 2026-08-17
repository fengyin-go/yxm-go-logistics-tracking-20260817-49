package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"logistics/internal/model"
	"logistics/internal/store"
)

// TestWriteServiceErrorStatusMapping 锁定错误到 HTTP 状态码的映射：
// ErrNotFound 必须映射为 404，而非 409，否则调用方无法识别“记录不存在”。
func TestWriteServiceErrorStatusMapping(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want int
	}{
		{"validation -> 400", model.NewValidationError("field", "invalid"), http.StatusBadRequest},
		{"not found -> 404", store.ErrNotFound, http.StatusNotFound},
		{"wrapped not found -> 404", errors.Join(store.ErrNotFound, errors.New("ctx")), http.StatusNotFound},
		{"conflict -> 409", store.ErrConflict, http.StatusConflict},
		{"unknown -> 500", errors.New("boom"), http.StatusInternalServerError},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			writeServiceError(rec, c.err)
			if rec.Code != c.want {
				t.Fatalf("status = %d, want %d (body=%q)", rec.Code, c.want, rec.Body.String())
			}
		})
	}
}
