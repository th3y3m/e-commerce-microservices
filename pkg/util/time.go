package util

import (
	"reflect"
	"time"
)

const tsCreateTimeLayout = "2006-01-02 15:04:05 +0700"

func FormatVNTime(t time.Time, layout string) string {
	if reflect.DeepEqual(t, time.Time{}) {
		return ""
	}
	t = t.Add(7 * time.Hour)
	return t.Format(layout)
}

func ParseTime(timeStr string) time.Time {
	parsedTime, err := time.Parse(tsCreateTimeLayout, timeStr)
	if err != nil {
		return time.Time{}
	}
	return parsedTime
}
