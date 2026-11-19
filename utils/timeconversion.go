// Copyright (c) 2025 DaggerTech. All rights reserved.
// Use of this source code is governed by an MIT license that can be
// found in the LICENSE file.
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
	nss := strconv.FormatInt(int64(ns), 10)
	if len(nss) > 3 {
		nss = nss[:3]
	}

	// Format the time components into SQL string
	return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d.%s", y, m, d, h, mn, s, nss)
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
		secStr := ""
		nsStr := ""
		inNS := false
		nums := "0123456789"
		for _, c := range timeParts[2] {
			if c == '.' && inNS {
				break
			}
			if c == '.' {
				inNS = true
				continue
			}
			if !strings.Contains(nums, string(c)) {
				break
			}
			if inNS {
				nsStr += string(c)
			} else {
				secStr += string(c)
			}
		}
		s, e = strconv.Atoi(secStr)
		if e != nil {
			return nil, false
		}

		ns = 0
		if nsStr != "" {
			ns, e = strconv.Atoi(nsStr)
			if e != nil {
				return nil, false
			}
		}
	}
	// Create and return time.Time object
	t := time.Date(y, time.Month(m), d, h, mn, s, ns, time.UTC)
	return &t, true
}
