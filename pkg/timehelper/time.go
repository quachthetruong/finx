package timehelper

import (
	"time"
)

func GetDurationInMilliseconds(start time.Time) int64 {
	end := time.Now()
	duration := end.Sub(start)
	return duration.Milliseconds()
}
