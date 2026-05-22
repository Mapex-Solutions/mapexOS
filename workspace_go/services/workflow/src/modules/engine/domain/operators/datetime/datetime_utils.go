package datetime

import (
	"time"
)

/**
 * Common date formats supported by the datetime operators.
 */
var DateFormats = []string{
	time.RFC3339,
	time.RFC3339Nano,
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05",
	"2006-01-02",
	"02/01/2006",
	"01/02/2006",
	"2006/01/02",
}

/**
 * Common time-only formats supported by time operators.
 */
var TimeFormats = []string{
	"15:04:05",
	"15:04",
	"3:04 PM",
	"3:04:05 PM",
}

/**
 * ParseDate attempts to parse a value as a date.
 * Supports time.Time, string (multiple formats), and Unix timestamps.
 *
 * @param v - Value to parse
 * @param timezone - IANA timezone for parsing (e.g., "America/Sao_Paulo")
 * @returns (time.Time, true) if successful, (zero, false) otherwise
 */
func ParseDate(v interface{}, timezone string) (time.Time, bool) {
	loc := getLocation(timezone)

	switch val := v.(type) {
	case time.Time:
		return val.In(loc), true
	case string:
		// Try common date formats
		for _, format := range DateFormats {
			if t, err := time.ParseInLocation(format, val, loc); err == nil {
				return t, true
			}
		}
	case int64:
		// Unix timestamp (seconds)
		return time.Unix(val, 0).In(loc), true
	case float64:
		// Unix timestamp with milliseconds
		sec := int64(val)
		nsec := int64((val - float64(sec)) * 1e9)
		return time.Unix(sec, nsec).In(loc), true
	case int:
		return time.Unix(int64(val), 0).In(loc), true
	}
	return time.Time{}, false
}

/**
 * ParseTime attempts to parse a value as a time-of-day.
 * Returns hours, minutes, seconds.
 *
 * @param v - Value to parse (string like "14:30" or "2:30 PM")
 * @param timezone - IANA timezone for parsing
 * @returns (time.Time with today's date and parsed time, true) if successful
 */
func ParseTime(v interface{}, timezone string) (time.Time, bool) {
	loc := getLocation(timezone)

	switch val := v.(type) {
	case time.Time:
		// Return as-is, just ensure correct location
		return val.In(loc), true
	case string:
		// Try time-only formats
		for _, format := range TimeFormats {
			if t, err := time.ParseInLocation(format, val, loc); err == nil {
				// Set date to today for comparison
				now := time.Now().In(loc)
				return time.Date(
					now.Year(), now.Month(), now.Day(),
					t.Hour(), t.Minute(), t.Second(), 0, loc,
				), true
			}
		}
		// Also try full date formats (extract time component)
		for _, format := range DateFormats {
			if t, err := time.ParseInLocation(format, val, loc); err == nil {
				// Extract just the time portion with today's date
				now := time.Now().In(loc)
				return time.Date(
					now.Year(), now.Month(), now.Day(),
					t.Hour(), t.Minute(), t.Second(), 0, loc,
				), true
			}
		}
	}
	return time.Time{}, false
}

/**
 * CompareDates compares two dates.
 * Returns: -1 (a before b), 0 (equal), 1 (a after b)
 */
func CompareDates(a, b time.Time) int {
	if a.Before(b) {
		return -1
	}
	if a.After(b) {
		return 1
	}
	return 0
}

/**
 * CompareTimes compares only the time portion of two times.
 * Ignores the date component.
 * Returns: -1 (a before b), 0 (equal), 1 (a after b)
 */
func CompareTimes(a, b time.Time) int {
	aSeconds := a.Hour()*3600 + a.Minute()*60 + a.Second()
	bSeconds := b.Hour()*3600 + b.Minute()*60 + b.Second()

	if aSeconds < bSeconds {
		return -1
	}
	if aSeconds > bSeconds {
		return 1
	}
	return 0
}

/**
 * DateOnly extracts just the date portion (zeroing time).
 */
func DateOnly(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

/**
 * getLocation returns a time.Location for the given timezone.
 * Falls back to UTC if timezone is invalid.
 */
func getLocation(timezone string) *time.Location {
	if timezone == "" {
		return time.UTC
	}
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return time.UTC
	}
	return loc
}

/**
 * GetCurrentDateInTimezone returns the current date in the specified timezone.
 */
func GetCurrentDateInTimezone(timezone string) time.Time {
	loc := getLocation(timezone)
	return time.Now().In(loc)
}

/**
 * GetCurrentTimeInTimezone returns the current time in the specified timezone.
 */
func GetCurrentTimeInTimezone(timezone string) time.Time {
	loc := getLocation(timezone)
	return time.Now().In(loc)
}
