package jutils

import (
	"fmt"
	"time"
)

func FriendlyTimestamp() string {
	currentTime := time.Now()
	return fmt.Sprintf("%d-%d-%d %d:%d:%d\n",
		currentTime.Year(),
		currentTime.Month(),
		currentTime.Day(),
		currentTime.Hour(),
		currentTime.Minute(),
		currentTime.Second())
}
