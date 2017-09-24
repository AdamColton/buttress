package gsql

import (
	"time"
)

var TimeFormat = "2006-01-02 15:04:05"

func TimeToString(t time.Time) string {
	return t.Format(TimeFormat)
}

func StringToTime(str string) time.Time {
	t, _ := time.Parse(TimeFormat, str)
	return t
}
