// Package scheduler computes edition windows and drives the newsletter's
// open/reminder/grace/sent state machine.
package scheduler

import (
	"errors"
	"time"
)

// NextWindow returns the next opens_at instant strictly after "after" for the
// given period/anchor configuration. The anchor weekday/day-of-month is
// evaluated against the calendar date in tz (so "Friday" means Friday in the
// group's local calendar); the resulting instant uses that calendar date with
// sendHourUTC as the hour, in UTC — DST shifts in tz are intentionally not
// reflected in the fire hour, which keeps the math simple and is an
// acceptable tradeoff for a friend-group newsletter.
//
// period is one of "weekly", "biweekly", "monthly", "quarterly".
// anchorWeekday is 0=Sunday..6=Saturday (used by weekly/biweekly).
// anchorDayOfMonth is 1-31, clamped to the last day of short months (used by
// monthly/quarterly).
// epochDate anchors biweekly parity and quarterly month-of-quarter; only its
// date (in tz) matters.
func NextWindow(period string, anchorWeekday, anchorDayOfMonth int, epochDate time.Time, sendHourUTC int, tz string, after time.Time) (time.Time, error) {
	loc, err := time.LoadLocation(tz)
	if err != nil || tz == "" {
		loc = time.UTC
	}

	match, err := matcherFor(period, anchorWeekday, anchorDayOfMonth, dateOnly(epochDate, loc))
	if err != nil {
		return time.Time{}, err
	}

	candidate := dateOnly(after, loc)
	for range 400 {
		if match(candidate) {
			fire := time.Date(candidate.Year(), candidate.Month(), candidate.Day(), sendHourUTC, 0, 0, 0, time.UTC)
			if fire.After(after) {
				return fire, nil
			}
		}
		candidate = candidate.AddDate(0, 0, 1)
	}
	return time.Time{}, errors.New("no matching window found within 400 days")
}

// dateOnly returns midnight of t's calendar date in loc.
func dateOnly(t time.Time, loc *time.Location) time.Time {
	lt := t.In(loc)
	y, m, d := lt.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, loc)
}

func daysInMonth(y int, m time.Month) int {
	return time.Date(y, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

// clampedDay returns the target day-of-month for y/m, clamping anchorDay to
// the last day of short months (e.g. 31 in February becomes 28 or 29).
func clampedDay(y int, m time.Month, anchorDay int) int {
	last := daysInMonth(y, m)
	if anchorDay > last {
		return last
	}
	if anchorDay < 1 {
		return 1
	}
	return anchorDay
}

func matcherFor(period string, anchorWeekday, anchorDayOfMonth int, epoch time.Time) (func(time.Time) bool, error) {
	switch period {
	case "weekly":
		return func(d time.Time) bool {
			return int(d.Weekday()) == anchorWeekday
		}, nil
	case "biweekly":
		return func(d time.Time) bool {
			if int(d.Weekday()) != anchorWeekday {
				return false
			}
			days := int(d.Sub(epoch).Hours() / 24)
			weekIndex := days / 7
			if days < 0 {
				// floor toward negative infinity for dates before epoch
				weekIndex = -((-days + 6) / 7)
			}
			return weekIndex%2 == 0
		}, nil
	case "monthly":
		return func(d time.Time) bool {
			return d.Day() == clampedDay(d.Year(), d.Month(), anchorDayOfMonth)
		}, nil
	case "quarterly":
		return func(d time.Time) bool {
			if d.Day() != clampedDay(d.Year(), d.Month(), anchorDayOfMonth) {
				return false
			}
			monthsSinceEpoch := (d.Year()-epoch.Year())*12 + int(d.Month()) - int(epoch.Month())
			return ((monthsSinceEpoch % 3) + 3) % 3 == 0
		}, nil
	default:
		return nil, errors.New("unknown schedule_period: " + period)
	}
}
