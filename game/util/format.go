package util

import "time"

func Format(dur time.Duration, format string) string {
	return time.Unix(0, 0).UTC().Add(time.Duration(dur)).Format(format)
}
