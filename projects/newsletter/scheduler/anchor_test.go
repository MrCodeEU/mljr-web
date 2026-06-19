package scheduler

import (
	"testing"
	"time"
)

func utc(y int, m time.Month, d, h int) time.Time {
	return time.Date(y, m, d, h, 0, 0, 0, time.UTC)
}

func TestNextWindowWeekly(t *testing.T) {
	// 2026-06-19 is a Friday. anchorWeekday=5 (Friday), sendHour=18.
	after := utc(2026, 6, 16, 0) // Tuesday
	got, err := NextWindow("weekly", 5, 0, time.Time{}, 18, "UTC", after)
	if err != nil {
		t.Fatal(err)
	}
	want := utc(2026, 6, 19, 18)
	if !got.Equal(want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestNextWindowWeeklySkipsPastFireTimeSameDay(t *testing.T) {
	// after is already past this Friday's fire hour -> jump to next Friday.
	after := utc(2026, 6, 19, 19) // Friday 19:00, fire hour is 18:00
	got, err := NextWindow("weekly", 5, 0, time.Time{}, 18, "UTC", after)
	if err != nil {
		t.Fatal(err)
	}
	want := utc(2026, 6, 26, 18)
	if !got.Equal(want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestNextWindowBiweeklyParity(t *testing.T) {
	epoch := utc(2026, 6, 5, 0) // a Friday, parity week 0
	// next Friday (2026-06-12) is parity week 1 (odd) -> should skip to 06-19
	after := utc(2026, 6, 6, 0)
	got, err := NextWindow("biweekly", 5, 0, epoch, 18, "UTC", after)
	if err != nil {
		t.Fatal(err)
	}
	want := utc(2026, 6, 19, 18)
	if !got.Equal(want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestNextWindowMonthly(t *testing.T) {
	after := utc(2026, 1, 10, 0)
	got, err := NextWindow("monthly", 0, 15, time.Time{}, 9, "UTC", after)
	if err != nil {
		t.Fatal(err)
	}
	want := utc(2026, 1, 15, 9)
	if !got.Equal(want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestNextWindowMonthlyClampsFebruary(t *testing.T) {
	// anchorDayOfMonth=31; Feb 2026 (non-leap) has 28 days -> clamp to 28.
	after := utc(2026, 1, 31, 23)
	got, err := NextWindow("monthly", 0, 31, time.Time{}, 9, "UTC", after)
	if err != nil {
		t.Fatal(err)
	}
	want := utc(2026, 2, 28, 9)
	if !got.Equal(want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestNextWindowMonthlyClampsLeapFebruary(t *testing.T) {
	after := utc(2028, 1, 31, 23) // 2028 is a leap year
	got, err := NextWindow("monthly", 0, 31, time.Time{}, 9, "UTC", after)
	if err != nil {
		t.Fatal(err)
	}
	want := utc(2028, 2, 29, 9)
	if !got.Equal(want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestNextWindowQuarterly(t *testing.T) {
	epoch := utc(2026, 1, 1, 0) // quarter months: Jan, Apr, Jul, Oct
	after := utc(2026, 2, 1, 0)
	got, err := NextWindow("quarterly", 0, 1, epoch, 12, "UTC", after)
	if err != nil {
		t.Fatal(err)
	}
	want := utc(2026, 4, 1, 12)
	if !got.Equal(want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestNextWindowUnknownPeriod(t *testing.T) {
	if _, err := NextWindow("daily", 0, 0, time.Time{}, 0, "UTC", time.Now()); err == nil {
		t.Error("expected error for unknown period")
	}
}
