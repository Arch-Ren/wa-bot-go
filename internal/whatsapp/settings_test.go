package whatsapp

import (
	"testing"
	"time"
)

func TestTodayKeyLocation(t *testing.T) {
	// GMT+7 = 2026-09-03 22:30 WIB
	ref := time.Date(2026, 9, 3, 15, 30, 0, 0, time.UTC)
	want := "2026-09-03"
	got := ref.Add(gmt7Offset).Format("2006-01-02")
	if got != want {
		t.Errorf("todayKey mismatch: got %s want %s", got, want)
	}
}

func TestTodayKeyBoundary(t *testing.T) {
	// UTC 2026-09-03 17:00 = GMT+7 2026-09-04 00:00 -> harus 2026-09-04
	ref := time.Date(2026, 9, 3, 17, 0, 0, 0, time.UTC)
	want := "2026-09-04"
	got := ref.Add(gmt7Offset).Format("2006-01-02")
	if got != want {
		t.Errorf("boundary mismatch: got %s want %s", got, want)
	}
}
