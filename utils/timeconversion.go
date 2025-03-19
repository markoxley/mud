package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func TimeToSQL(t time.Time) string {
	var y, m, d, h, mn, s, ns int
	y = t.Year()
	m = int(t.Month())
	d = t.Day()
	h = t.Hour()
	mn = t.Minute()
	s = t.Second()
	ns = t.Nanosecond()
	return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d.%d", y, m, d, h, mn, s, ns)
}

func SQLToTime(st string) (*time.Time, bool) {
	sep := " "
	var y, m, d, h, mn, s, ns int
	var e error
	if strings.Contains(st, "T") {
		sep = "T"
	}
	dt := strings.Split(st, sep)
	if len(dt) < 1 {
		return nil, false
	}
	dp := strings.Split(dt[0], "-")
	if len(dp) != 3 {
		return nil, false
	}
	if y, e = strconv.Atoi(dp[0]); e != nil {
		return nil, false
	}
	if m, e = strconv.Atoi(dp[1]); e != nil {
		return nil, false
	}
	if d, e = strconv.Atoi(dp[2]); e != nil {
		return nil, false
	}
	if len(dt) > 1 {
		tm := strings.Split(dt[1], ":")
		if len(tm) != 3 {
			return nil, false
		}

		sc := strings.Split(tm[2], ".")

		if len(sc) > 1 {
			if ns, e = strconv.Atoi(sc[1]); e != nil {
				ns = 0
			}
		}
		if h, e = strconv.Atoi(tm[0]); e != nil {
			return nil, false
		}
		if mn, e = strconv.Atoi(tm[1]); e != nil {
			return nil, false
		}
		if s, e = strconv.Atoi(sc[0]); e != nil {
			return nil, false
		}
	}
	t := time.Date(y, time.Month(m), d, h, mn, s, ns, time.UTC)
	return &t, true
}
