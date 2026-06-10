package handlers

import (
	"fmt"
	"time"
)

const maxRangeMonths = 24

// parseMonthRange parses "YYYY-MM" query params; empty values default to the
// last 12 months including the current one. Returns [from, toExclusive).
func parseMonthRange(fromStr, toStr string, now time.Time) (time.Time, time.Time, error) {
	to := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	from := to.AddDate(0, -11, 0)
	var err error
	if fromStr != "" {
		if from, err = time.Parse("2006-01", fromStr); err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid from: %v", err)
		}
	}
	if toStr != "" {
		if to, err = time.Parse("2006-01", toStr); err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid to: %v", err)
		}
	}
	if to.Before(from) {
		return time.Time{}, time.Time{}, fmt.Errorf("from is after to")
	}
	months := (to.Year()-from.Year())*12 + int(to.Month()) - int(from.Month()) + 1
	if months > maxRangeMonths {
		return time.Time{}, time.Time{}, fmt.Errorf("range exceeds %d months", maxRangeMonths)
	}
	return from, to.AddDate(0, 1, 0), nil
}
