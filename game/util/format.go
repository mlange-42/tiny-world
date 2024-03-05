package util

import (
	"fmt"
	"strings"
	"time"
)

func FormatDuration(dur time.Duration) string {
	hours := int(dur.Hours())
	mins := int(dur.Minutes()) - 60*hours
	return fmt.Sprintf("%d:%02d", hours, mins)
}

func Capitalize(s string) string {
	if len(s) == 0 {
		return ""
	}
	runes := []rune(s)
	runes[0] = []rune(strings.ToUpper(string(runes[0])))[0]
	return string(runes)
}
