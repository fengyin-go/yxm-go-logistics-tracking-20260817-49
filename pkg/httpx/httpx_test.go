package httpx

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestParsePagination(t *testing.T) {
	cases := []struct {
		name                 string
		query                string
		defaultSize, maxSize int
		wantPage, wantSize   int
	}{
		{"empty uses defaults", "", 20, 100, 1, 20},
		{"page=0 size=0 clamps", "page=0&size=0", 20, 100, 1, 20},
		{"negative clamps", "page=-1&size=-5", 20, 100, 1, 20},
		{"normal values", "page=2&size=50", 20, 100, 2, 50},
		{"size capped to max", "page=3&size=500", 20, 100, 3, 100},
		{"non-numeric falls back", "page=abc&size=xy", 20, 100, 1, 20},
		{"invalid defaultSize falls back", "page=1&size=0", 0, 100, 1, 20},
		{"invalid maxSize falls back", "page=1&size=500", 20, 0, 1, 100},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/?"+c.query, nil)
			pp := ParsePagination(req, c.defaultSize, c.maxSize)
			if pp.Page != c.wantPage || pp.Size != c.wantSize {
				t.Errorf("got page=%d size=%d, want page=%d size=%d",
					pp.Page, pp.Size, c.wantPage, c.wantSize)
			}
			if pp.Page < 1 || pp.Size < 1 {
				t.Errorf("invalid result: page=%d size=%d", pp.Page, pp.Size)
			}
		})
	}
}
