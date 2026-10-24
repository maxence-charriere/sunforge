package vi

import (
	"fmt"
	"time"
)

func TimeSince(t time.Time) string {
	d := time.Since(t)
	d = d.Round(d)

	var unit string
	var count float64

	switch {
	case d.Hours() > 24*30*12:
		count = d.Hours() / float64(24*30*12)
		unit = "year"

	case d.Hours() > 24*30:
		count = d.Hours() / float64(24*30)
		unit = "month"

	case d.Hours() >= 24:
		count = d.Hours() / float64(24)
		unit = "day"

	case d.Hours() >= 1:
		count = d.Hours()
		unit = "hour"

	case d.Minutes() >= 1:
		count = d.Minutes()
		unit = "minute"

	case d.Seconds() >= 1:
		count = d.Seconds()
		unit = "second"

	default:
		count = 1
		unit = "second"
	}

	if count >= 2 {
		unit += "s"
	}

	return fmt.Sprintf("%.f %s", count, unit)
}

func IfTrue[T any](b bool, v T) T {
	if b {
		return v
	}
	var zero T
	return zero
}
