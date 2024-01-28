package util

import (
	"strings"
	"time"
)

func GetTimeStamp() string {
	return strings.Replace(time.Now().Format("2006-01-02T15:04:05.999999999"), "Z", "", -1)
}
