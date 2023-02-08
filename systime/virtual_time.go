package systime

import (
	"time"
)

type VirtualTime interface {
	Now() time.Time
	Time(ts int64) time.Time
	Parse(string) (time.Time, error)
	Format(time.Time) string
	Location() *time.Location
	FakeDuration() time.Duration
}

const (
	datetimeParse  = "2006-1-2 15:04:05"
	dateParse      = "2006-1-2"
	datetimeFormat = "2006-01-02 15:04:05"
)

var virtualTime VirtualTime = NewFakeTime()

func SetVirtualTime(vt VirtualTime) {
	virtualTime = vt
}

type defaultVirtualTime struct{}

func (defaultVirtualTime) Parse(s string) (time.Time, error) {
	var format string
	if len(s) > 10 {
		format = datetimeParse
	} else {
		format = dateParse
	}
	return time.Parse(format, s)
}

func (defaultVirtualTime) Format(tm time.Time) string {
	return tm.Format(datetimeFormat)
}

func (defaultVirtualTime) Now() time.Time {
	return time.Now()
}

func (defaultVirtualTime) Time(ts int64) time.Time {
	return time.Unix(ts, 0)
}

func (defaultVirtualTime) Location() *time.Location {
	return time.Local
}

func (defaultVirtualTime) FakeDuration() time.Duration {
	return 0
}
