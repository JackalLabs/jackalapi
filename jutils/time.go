package jutils

import (
	"fmt"
	"time"
)

func FriendlyTimestamp() string {
	currentTime := time.Now()
	return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d",
		currentTime.Year(),
		currentTime.Month(),
		currentTime.Day(),
		currentTime.Hour(),
		currentTime.Minute(),
		currentTime.Second())
}

func UnixMsTimestamp() int64 {
	currentTime := time.Now()
	return currentTime.UnixMilli()
}
