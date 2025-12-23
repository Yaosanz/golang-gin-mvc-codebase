package helpers

import (
	"log"
	"time"
)

// ParseTime parses a string into time.Time.
// - Prioritizes RFC3339 (default JSON format).
// - Falls back to other common formats.
// - Normalizes to UTC for safe storage in Postgres timestamptz.
func ParseTime(value string) time.Time {
	if value == "" {
		log.Print("empty time string")
		return time.Time{}
	}

	formats := []string{
		time.RFC3339,          // e.g. 2025-09-10T14:30:00Z
		"2006-01-02 15:04:05", // e.g. 2025-09-10 14:30:00
		"2006-01-02",          // e.g. 2025-09-10
		"02-01-2006",          // e.g. 10-09-2025
		"02/01/2006",          // e.g. 10/09/2025
		"2006/01/02",          // e.g. 2025/09/10
		"15:04:05",            // e.g. 14:30:00
	}

	var t time.Time
	var err error
	for _, f := range formats {
		t, err = time.ParseInLocation(f, value, time.Local)
		if err == nil {
			return t.UTC() // normalize for timestamptz
		}
	}

	log.Print("unsupported time format: " + value)
	return time.Time{}
}
