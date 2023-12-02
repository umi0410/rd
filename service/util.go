package service

import "time"

func SetWhenZeroTimeValue(dst *time.Time, val time.Time) {
	if dst.IsZero() {
		*dst = val
	}
}
