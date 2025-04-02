// Package utils provides utility functions for time conversion between SQL and Go
package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// TimeToSQL converts a Go time.Time object to a SQL-compatible string format
// The format is "YYYY-MM-DD HH:MM:SS.NNNNNNNNN" where N is nanoseconds
//
// @param t The time.Time object to convert
// @return A SQL-compatible string representation of the time
func TimeToSQL(t time.Time) string {
	var y, m, d, h, mn, s, ns int
	// Extract time components
	y = t.Year()
	m = int(t.Month())
	d = t.Day()
	h = t.Hour()
	mn = t.Minute()
	s = t.Second()
	ns = t.Nanosecond()

	// Format the time components into SQL string
	return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d.%d", y, m, d, h, mn, s, ns)
}

// SQLToTime converts a SQL time string to a Go time.Time object
// Supports both " " and "T" as separators between date and time
// Returns a pointer to time.Time and a boolean indicating success
//
// @param st The SQL time string to convert
// @return A pointer to the converted time.Time object, and a boolean indicating success
func SQLToTime(st string) (*time.Time, bool) {
	sep := " "
	var y, m, d, h, mn, s, ns int
	var e error

	// Check for alternative separator
	if strings.Contains(st, "T") {
		sep = "T"
	}

	// Split into date and time parts
	dt := strings.Split(st, sep)
	if len(dt) < 1 {
		return nil, false
	}

	// Parse date part
	if len(dt) > 0 {
		dateParts := strings.Split(dt[0], "-")
		if len(dateParts) != 3 {
			return nil, false
		}
		y, e = strconv.Atoi(dateParts[0])
		if e != nil {
			return nil, false
		}
		m, e = strconv.Atoi(dateParts[1])
		if e != nil {
			return nil, false
		}
		d, e = strconv.Atoi(dateParts[2])
		if e != nil {
			return nil, false
		}
	}

	// Parse time part
	if len(dt) > 1 {
		timeParts := strings.Split(dt[1], ":")
		if len(timeParts) != 3 {
			return nil, false
		}
		h, e = strconv.Atoi(timeParts[0])
		if e != nil {
			return nil, false
		}
		mn, e = strconv.Atoi(timeParts[1])
		if e != nil {
			return nil, false
		}
		s, e = strconv.Atoi(timeParts[2])
		if e != nil {
			return nil, false
		}
		if strings.Contains(timeParts[2], ".") {
			nsParts := strings.Split(timeParts[2], ".")
			if len(nsParts) > 1 {
				ns, e = strconv.Atoi(nsParts[1])
				if e != nil {
					return nil, false
				}
			}
		}
	}

	// Create and return time.Time object
	t := time.Date(y, time.Month(m), d, h, mn, s, ns, time.UTC)
	return &t, true
}
