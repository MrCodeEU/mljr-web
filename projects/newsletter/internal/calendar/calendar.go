// Package calendar holds pure calendar-date arithmetic shared by the
// newsletter's profile/edition features (starsign inference, upcoming-
// birthday windows) — no PocketBase dependency, so it's trivially testable
// in isolation.
package calendar

import "time"

// SignForDate returns the Western zodiac sign for a given month (1-12) and
// day, using fixed calendar-date boundaries. Returns "" for an invalid
// month/day.
func SignForDate(month, day int) string {
	switch {
	case month < 1 || month > 12 || day < 1 || day > 31:
		return ""
	case (month == 3 && day >= 21) || (month == 4 && day <= 19):
		return "Aries"
	case (month == 4 && day >= 20) || (month == 5 && day <= 20):
		return "Taurus"
	case (month == 5 && day >= 21) || (month == 6 && day <= 20):
		return "Gemini"
	case (month == 6 && day >= 21) || (month == 7 && day <= 22):
		return "Cancer"
	case (month == 7 && day >= 23) || (month == 8 && day <= 22):
		return "Leo"
	case (month == 8 && day >= 23) || (month == 9 && day <= 22):
		return "Virgo"
	case (month == 9 && day >= 23) || (month == 10 && day <= 22):
		return "Libra"
	case (month == 10 && day >= 23) || (month == 11 && day <= 21):
		return "Scorpio"
	case (month == 11 && day >= 22) || (month == 12 && day <= 21):
		return "Sagittarius"
	case (month == 12 && day >= 22) || (month == 1 && day <= 19):
		return "Capricorn"
	case (month == 1 && day >= 20) || (month == 2 && day <= 18):
		return "Aquarius"
	case (month == 2 && day >= 19) || (month == 3 && day <= 20):
		return "Pisces"
	default:
		return ""
	}
}

// InRange reports whether the calendar date (month, day) falls within the
// [from, to) window, comparing only month/day (birthdays are stored with a
// placeholder year, age isn't tracked) and correctly handling windows that
// wrap across a Dec 31 -> Jan 1 boundary.
func InRange(month, day int, from, to time.Time) bool {
	if !from.Before(to) {
		return false
	}
	// Normalize everything onto a fixed reference year so day-of-year
	// comparisons are well-defined, then handle wraparound by comparing
	// the (possibly wrapped) span against the candidate date placed in
	// both the from-year and the from-year+1.
	const refYear = 2001 // arbitrary non-leap year
	candidate := time.Date(refYear, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	fromNorm := time.Date(refYear, from.Month(), from.Day(), 0, 0, 0, 0, time.UTC)
	toNorm := time.Date(refYear, to.Month(), to.Day(), 0, 0, 0, 0, time.UTC)
	if !toNorm.After(fromNorm) {
		toNorm = toNorm.AddDate(1, 0, 0)
	}

	if !candidate.Before(fromNorm) && candidate.Before(toNorm) {
		return true
	}
	// Also check the candidate shifted a year later, to catch windows that
	// wrap past Dec 31 while the candidate date is early in the year.
	candidateNext := candidate.AddDate(1, 0, 0)
	return !candidateNext.Before(fromNorm) && candidateNext.Before(toNorm)
}
