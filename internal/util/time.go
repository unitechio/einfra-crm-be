
package util

import (
	"time"
)

// GetCurrentUTC returns the current time in UTC.
func GetCurrentUTC() time.Time {
	return time.Now().UTC()
}

// ParseTime parses a string into a time.Time object using the RFC3339 layout.
func ParseTime(timeStr string) (time.Time, error) {
	return time.Parse(time.RFC3339, timeStr)
}

// FormatTime formats a time.Time object into a string using the RFC3339 layout.
func FormatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}
