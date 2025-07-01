package hntime

import (
	"time"
)

var loc, _ = time.LoadLocation("Asia/Bangkok")

const clearRequestTime = 15

func Now() time.Time {
	return time.Now().In(loc)
}

func TimeClearRequest() time.Time {
	currentHnTime := Now()
	if hour := currentHnTime.Hour(); hour > clearRequestTime {
		return time.Date(
			currentHnTime.Year(), currentHnTime.Month(), currentHnTime.Day(), clearRequestTime, 0, 0, 0, loc,
		).UTC()
	}
	yesterday := currentHnTime.Add(-24 * time.Hour)
	return time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), clearRequestTime, 0, 0, 0, loc).UTC()
}
