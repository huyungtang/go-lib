package times

import (
	"fmt"
	"strconv"
	base "time"
)

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

var (
	RCF3339  = base.RFC3339
	RCF3339S = "2006/01/02T15:04:05Z07:00"
	DateFmt  = "2006-01-02 15:04:05"
	DateFmtS = "2006/01/02 15:04:05"
	DateStr  = "20060102150405"
)

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// Add
// ****************************************************************************************************************************************
func Add(tm base.Time, ymdhms ...int64) base.Time {
	y, m, d := 0, 0, 0
	dur := base.Duration(0)
	for k, v := range ymdhms {
		switch k {
		case 0:
			y = int(v)
		case 1:
			m = int(v)
		case 2:
			d = int(v)
		case 3:
			dur = dur + (base.Hour * base.Duration(v))
		case 4:
			dur = dur + (base.Minute * base.Duration(v))
		case 5:
			dur = dur + (base.Second * base.Duration(v))
		}
	}

	return tm.AddDate(y, m, d).Add(dur)
}

// AsLocal
// ****************************************************************************************************************************************
func AsLocal(t base.Time) base.Time {
	return base.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), base.Local)
}

// Duration
// ****************************************************************************************************************************************
func Duration(h, m, s int64) (d base.Duration) {
	d, _ = base.ParseDuration(fmt.Sprintf("%dh%dm%ds", h, m, s))

	return
}

// Format
// ****************************************************************************************************************************************
func Format(tm base.Time, layout string) string {
	return tm.Format(layout)
}

// FromUnix
// ****************************************************************************************************************************************
func FromUnix(second, nanoSecond int64) base.Time {
	return base.Unix(second, nanoSecond)
}

// FromUnixDay
// ****************************************************************************************************************************************
func FromUnixDay(day int64) base.Time {
	return base.Unix(0, day*int64(base.Hour)*24)
}

// Now
// ****************************************************************************************************************************************
func Now() base.Time {
	return base.Now()
}

// Parse
// ****************************************************************************************************************************************
func Parse(layout, value string) (v base.Time, e error) {
	return base.Parse(layout, value)
}

// Period
// ****************************************************************************************************************************************
func Period(tm base.Time) (v uint64) {
	dt := Format(tm, DateStr[0:6])
	v, _ = strconv.ParseUint(dt, 10, 64)

	return
}

// Today
// ****************************************************************************************************************************************
func Today() base.Time {
	y, m, d := base.Now().Date()

	return base.Date(y, m, d, 0, 0, 0, 0, base.Local)
}

// UnixDay
// ****************************************************************************************************************************************
func UnixDay(tm base.Time) int64 {
	y, m, d := tm.Date()

	return base.Date(y, m, d, 0, 0, 0, 0, base.UTC).UnixNano() / (int64(base.Hour) * 24)
}

// UnixSecond
// ****************************************************************************************************************************************
func UnixSecond(tm base.Time) int64 {
	return tm.Unix()
}

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
