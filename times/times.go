package times

import (
	"fmt"
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

// Duration
// ****************************************************************************************************************************************
func Duration(h, m, s int64) (d base.Duration) {
	d, _ = base.ParseDuration(fmt.Sprintf("%dh%dm%ds", h, m, s))

	return
}

// Now
// ****************************************************************************************************************************************
func Now() Time {
	return time(base.Now())
}

// Today
// ****************************************************************************************************************************************
func Today() Time {
	return time(base.Now().Truncate(base.Hour * 24))
}

// Unix
// ****************************************************************************************************************************************
func Unix(second, nanoSecond int64) Time {
	return time(base.Unix(second, nanoSecond))
}

// Local
// ****************************************************************************************************************************************
func Local(t base.Time) Time {
	return time(base.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), base.Local))
}

// Parse
// ****************************************************************************************************************************************
func Parse(layout, value string) (Time, error) {
	v, e := base.Parse(layout, value)
	if e != nil {
		return nil, e
	}

	return time(v), nil
}

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// Time
// ****************************************************************************************************************************************
type Time interface {
	Add(...int64) Time
	Local() Time
	Time() base.Time
	Format(string) string
	UnixDay() int64
	UnixSecond() int64
	UnixMilli() int64
	UnixMicro() int64
	UnixNano() int64
}

// time ***********************************************************************************************************************************
type time base.Time

// Add
// ****************************************************************************************************************************************
func (o time) Add(ymdhms ...int64) Time {
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

	return time(base.Time(o).AddDate(y, m, d).Add(dur))
}

// Local
// ****************************************************************************************************************************************
func (o time) Local() Time {

	return Local(base.Time(o))
}

// Time
// ****************************************************************************************************************************************
func (o time) Time() base.Time {
	return base.Time(o)
}

// Format
// ****************************************************************************************************************************************
func (o time) Format(fmt string) string {
	return base.Time(o).Format(fmt)
}

// UnixDay
// ****************************************************************************************************************************************
func (o time) UnixDay() int64 {
	return base.Time(o).UnixNano() / (int64(base.Hour) * 24)
}

// UnixSecond
// ****************************************************************************************************************************************
func (o time) UnixSecond() int64 {
	return base.Time(o).Unix()
}

// UnixMilli
// ****************************************************************************************************************************************
func (o time) UnixMilli() int64 {
	return base.Time(o).UnixMilli()
}

// UnixMicro
// ****************************************************************************************************************************************
func (o time) UnixMicro() int64 {
	return base.Time(o).UnixMicro()
}

// UnixNano
// ****************************************************************************************************************************************
func (o time) UnixNano() int64 {
	return base.Time(o).UnixNano()
}

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
