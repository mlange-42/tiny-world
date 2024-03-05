package util

import (
	"fmt"
	"time"
)

func FormatDuration(dur time.Duration) string {
	hours := int(dur.Hours())
	mins := int(dur.Minutes()) - 60*hours
	return fmt.Sprintf("%d:%02d", hours, mins)
}
