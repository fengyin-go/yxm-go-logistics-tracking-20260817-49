package service

import (
	"math"
	"testing"
)

func TestPageBounds(t *testing.T) {
	cases := []struct {
		name               string
		total, page, size  int
		wantStart, wantEnd int
	}{
		{"first page", 10, 1, 5, 0, 5},
		{"second page", 10, 2, 5, 5, 10},
		{"last page partial", 7, 2, 5, 5, 7},
		{"page beyond range returns empty", 10, 3, 5, 10, 10},
		{"page=0 clamps to 1", 3, 0, 2, 0, 2},
		{"negative page clamps to 1", 3, -5, 2, 0, 2},
		{"size=0 falls back to default", 3, 1, 0, 0, 3},
		{"negative size falls back to default", 3, 1, -10, 0, 3},
		{"empty total", 0, 1, 20, 0, 0},
		{"empty total with bad page/size", 0, 0, 0, 0, 0},
		{"huge page does not overflow panic", 5, math.MaxInt32, 100, 5, 5},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			start, end := pageBounds(c.total, c.page, c.size)
			if start != c.wantStart || end != c.wantEnd {
				t.Errorf("pageBounds(%d,%d,%d) = [%d,%d), want [%d,%d)",
					c.total, c.page, c.size, start, end, c.wantStart, c.wantEnd)
			}
			if start < 0 || end < start || end > c.total {
				t.Errorf("invariant violated: start=%d end=%d total=%d", start, end, c.total)
			}
		})
	}
}
