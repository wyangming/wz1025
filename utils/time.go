package utils

import (
	"fmt"
	"time"
	"strings"
)

//把字符串转化为时间，根据字符串来精确，最精确的格式为"YYYY-MM-DD HH:MM:SS.MMMMMM"
func FormatDataTime(time_str string) time.Time {
	time_str = strings.TrimSpace(time_str)
	var (
		t   time.Time
		err error
	)
	timeFormat := "2006-01-02 15:04:05.9999999"

	switch len(time_str) {
	case 10, 19, 21, 22, 23, 24, 25, 26: //up to "YYYY-MM-DD HH:MM:SS.MMMMMM"
		t, err = time.ParseInLocation(timeFormat[:len(time_str)], time_str, time.Local)
	default:
		err = fmt.Errorf("invalid time string: %s", time_str)
	}
	if err != nil {
		ErrorLog("time.go FormatDataTime method", err)
		return t
	}
	return t
}
