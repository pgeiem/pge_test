package engine

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

// Duration represents a duration with second resolution.
type Duration time.Duration

// toDuration converts Duration to time.Duration.
func (d Duration) toDuration() time.Duration {
	return time.Duration(d)
}

var durationRegex = regexp.MustCompile(`^(\d+w)?(\d+d)?(\d+h)?(\d+m)?(\d+s)?$`)

// ParseDuration parses a duration string with units: seconds (s), minutes (m), hours (h), days (d), weeks (w).
func ParseDuration(s string) (Duration, error) {
	var totalSeconds int64
	var multipliers = map[byte]int64{
		'w': 7 * 24 * 60 * 60,
		'd': 24 * 60 * 60,
		'h': 60 * 60,
		'm': 60,
		's': 1,
	}

	matches := durationRegex.FindStringSubmatch(s)
	if matches == nil || s == "" {
		return 0, fmt.Errorf("error while parsing %s duration, invalid pattern", s)
	}

	for _, match := range matches[1:] {
		if match == "" {
			continue
		}

		num, err := strconv.ParseInt(match[:len(match)-1], 10, 64)
		if err != nil {
			return 0, fmt.Errorf("error while parsing %s duration, invalid number, %s", s, err)
		}

		unit := match[len(match)-1]
		multiplier, exists := multipliers[unit]
		if !exists {
			return 0, fmt.Errorf("error while parsing %s duration, invalid unit (%c)", s, unit)
		}

		totalSeconds += num * multiplier
	}

	return Duration(totalSeconds * int64(time.Second)), nil
}
