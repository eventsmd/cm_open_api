package handlers

import (
	"testing"
	"time"
)

func TestParseMonthRange(t *testing.T) {
	now := time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC)

	t.Run("defaults to last 12 months", func(t *testing.T) {
		from, toExcl, err := parseMonthRange("", "", now)
		if err != nil {
			t.Fatal(err)
		}
		if got := from.Format("2006-01-02"); got != "2025-07-01" {
			t.Errorf("from = %s, want 2025-07-01", got)
		}
		if got := toExcl.Format("2006-01-02"); got != "2026-07-01" {
			t.Errorf("toExcl = %s, want 2026-07-01", got)
		}
	})

	t.Run("explicit range, exclusive upper bound", func(t *testing.T) {
		from, toExcl, err := parseMonthRange("2025-07", "2026-01", now)
		if err != nil {
			t.Fatal(err)
		}
		if got := from.Format("2006-01-02"); got != "2025-07-01" {
			t.Errorf("from = %s", got)
		}
		if got := toExcl.Format("2006-01-02"); got != "2026-02-01" {
			t.Errorf("toExcl = %s", got)
		}
	})

	for _, bad := range []struct{ from, to, name string }{
		{"2026-13", "", "bad month"},
		{"garbage", "", "garbage"},
		{"2026-05", "2026-01", "from after to"},
		{"2020-01", "2026-01", "span over 24 months"},
	} {
		t.Run("rejects "+bad.name, func(t *testing.T) {
			if _, _, err := parseMonthRange(bad.from, bad.to, now); err == nil {
				t.Error("expected error, got nil")
			}
		})
	}
}
