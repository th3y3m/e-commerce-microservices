package util

import (
	"reflect"
	"time"
)

func FormatVNTime(t time.Time, layout string) string {
	if reflect.DeepEqual(t, time.Time{}) {
		return ""
	}
	t = t.Add(7 * time.Hour)
	return t.Format(layout)
}
