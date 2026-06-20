package calendar

import (
	"testing"
	"time"
)

func TestSignForDate(t *testing.T) {
	cases := []struct {
		month, day int
		want       string
	}{
		{3, 21, "Aries"}, {4, 19, "Aries"},
		{4, 20, "Taurus"}, {5, 20, "Taurus"},
		{5, 21, "Gemini"}, {6, 20, "Gemini"},
		{6, 21, "Cancer"}, {7, 22, "Cancer"},
		{7, 23, "Leo"}, {8, 22, "Leo"},
		{8, 23, "Virgo"}, {9, 22, "Virgo"},
		{9, 23, "Libra"}, {10, 22, "Libra"},
		{10, 23, "Scorpio"}, {11, 21, "Scorpio"},
		{11, 22, "Sagittarius"}, {12, 21, "Sagittarius"},
		{12, 22, "Capricorn"}, {1, 19, "Capricorn"},
		{1, 20, "Aquarius"}, {2, 18, "Aquarius"},
		{2, 19, "Pisces"}, {3, 20, "Pisces"},
		{0, 1, ""}, {13, 1, ""}, {1, 0, ""}, {1, 32, ""},
	}
	for _, c := range cases {
		if got := SignForDate(c.month, c.day); got != c.want {
			t.Errorf("SignForDate(%d, %d) = %q, want %q", c.month, c.day, got, c.want)
		}
	}
}

func TestInRange(t *testing.T) {
	d := func(month, day int) time.Time { return time.Date(2026, time.Month(month), day, 0, 0, 0, 0, time.UTC) }
	dYear := func(year, month, day int) time.Time { return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC) }

	cases := []struct {
		name       string
		month, day int
		from, to   time.Time
		want       bool
	}{
		{"inside simple window", 6, 15, d(6, 1), d(6, 30), true},
		{"before window", 5, 1, d(6, 1), d(6, 30), false},
		{"after window", 7, 1, d(6, 1), d(6, 30), false},
		{"at from boundary (inclusive)", 6, 1, d(6, 1), d(6, 30), true},
		{"at to boundary (exclusive)", 6, 30, d(6, 1), d(6, 30), false},
		{"wraps year, candidate in december", 12, 28, d(12, 20), dYear(2027, 1, 5), true},
		{"wraps year, candidate in january", 1, 2, d(12, 20), dYear(2027, 1, 5), true},
		{"wraps year, candidate outside", 6, 1, d(12, 20), dYear(2027, 1, 5), false},
		{"invalid window (from after to)", 6, 15, d(6, 30), d(6, 1), false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := InRange(c.month, c.day, c.from, c.to); got != c.want {
				t.Errorf("InRange(%d, %d, %v, %v) = %v, want %v", c.month, c.day, c.from, c.to, got, c.want)
			}
		})
	}
}
