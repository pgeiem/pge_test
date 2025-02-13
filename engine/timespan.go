// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//Copied from https://github.com/rickb777/date/blob/v2/timespan/timespan.go

package engine

import (
	"fmt"
	"strings"
	"time"
)

// TimestampFormat is a simple format for date & time, "2006-01-02 15:04:05".
const TimestampFormat = "2006-01-02 15:04:05"

//const ISOFormat = "2006-01-02T15:04:05"

// TimeSpan holds a span of time between two instants with a 1 nanosecond resolution.
// It is implemented using a time.Duration, therefore is limited to a maximum span of 292 years.
//
// It supports RFC5545 timespan representations.
type TimeSpan struct {
	mark     time.Time
	duration time.Duration
}

// ZeroTimeSpan creates a new zero-duration time span at a specified time.
func ZeroTimeSpan(start time.Time) TimeSpan {
	return TimeSpan{start, 0}
}

// TimeSpanOf creates a new time span at a specified time and duration. The duration can
// be negative, e.g. for an alarm event before the mark time.
func TimeSpanOf(start time.Time, d time.Duration) TimeSpan {
	return TimeSpan{start, d}
}

// BetweenTimes creates a new time span from two times. The start and end can be in either
// order; the result will be normalised. The inputs are half-open: the start is included and
// the end is excluded.
func BetweenTimes(t1, t2 time.Time) TimeSpan {
	if t2.Before(t1) {
		return TimeSpan{t2, t1.Sub(t2)}
	}
	return TimeSpan{t1, t2.Sub(t1)}
}

// Start gets the start time of the time span.
func (ts TimeSpan) Start() time.Time {
	if ts.duration < 0 {
		return ts.mark.Add(ts.duration)
	}
	return ts.mark
}

// End gets the end time of the time span. Strictly, this is one nanosecond after the
// range of time included in the time span; this implements the half-open model.
func (ts TimeSpan) End() time.Time {
	if ts.duration < 0 {
		return ts.mark
	}
	return ts.mark.Add(ts.duration)
}

// Mark gets the time marked by this timespan. Typically this is the same as Start, but
// it's the same as End for time spans with negative duration.
func (ts TimeSpan) Mark() time.Time {
	return ts.mark
}

// Duration gets the duration of the time span.
func (ts TimeSpan) Duration() time.Duration {
	return ts.duration
}

// IsEmpty returns true if this is an empty time span (zero duration).
func (ts TimeSpan) IsEmpty() bool {
	return ts.duration == 0
}

// Normalise ensures that the mark time is at the start time and the duration is positive.
// The normalised time span is returned.
func (ts TimeSpan) Normalise() TimeSpan {
	if ts.duration < 0 {
		return TimeSpan{ts.mark.Add(ts.duration), -ts.duration}
	}
	return ts
}

// ShiftBy moves the time span by moving both the start and end times similarly.
// A negative parameter is allowed.
func (ts TimeSpan) ShiftBy(d time.Duration) TimeSpan {
	return TimeSpan{ts.mark.Add(d), ts.duration}
}

// ExtendBy lengthens the time span by a specified amount. The parameter may be negative,
// in which case it is possible that the end of the time span will appear to be before the
// start. However, the result is normalised so that the resulting start is the lesser value.
func (ts TimeSpan) ExtendBy(d time.Duration) TimeSpan {
	return TimeSpan{ts.mark, ts.duration + d}.Normalise()
}

// ExtendWithoutWrapping lengthens the time span by a specified amount. The parameter may be
// negative, but if its magnitude is large than the time span's duration, it will be truncated
// so that the result has zero duration in that case. The start time is never altered.
func (ts TimeSpan) ExtendWithoutWrapping(d time.Duration) TimeSpan {
	tsn := ts.Normalise()
	if d < 0 && -d > tsn.duration {
		return TimeSpan{tsn.mark, 0}
	}
	return TimeSpan{tsn.mark, tsn.duration + d}
}

// String produces a human-readable description of a time span.
func (ts TimeSpan) String() string {
	return fmt.Sprintf("%s from %s to %s", ts.duration, ts.mark.Format(TimestampFormat), ts.End().Format(TimestampFormat))
}

// In returns a TimeSpan adjusted from its current location to a new location. Because
// location is considered to be a presentational attribute, the actual time itself is not
// altered by this function. This matches the behaviour of time.Time.In(loc).
func (ts TimeSpan) In(loc *time.Location) TimeSpan {
	t := ts.mark.In(loc)
	return TimeSpan{t, ts.duration}
}

// Contains tests whether a given moment of time is enclosed within the time span. The
// start time is inclusive; the end time is exclusive.
// If t has a different locality to the time-span, it is adjusted accordingly.
func (ts TimeSpan) Contains(t time.Time) bool {
	tl := t.In(ts.mark.Location())
	return ts.mark.Equal(tl) || ts.mark.Before(tl) && ts.End().After(tl)
}

// Merge combines two time spans by calculating a time span that just encompasses them both.
// As a special case, if one span is entirely contained within the other span, the larger of
// the two is returned. Otherwise, the result is the start of the earlier one to the end of the
// later one, even if the two spans don't overlap.
func (ts TimeSpan) Merge(other TimeSpan) TimeSpan {
	if ts.mark.After(other.mark) {
		// swap the ranges to simplify the logic
		return other.Merge(ts)

	} else if ts.End().After(other.End()) {
		// other is a proper subrange of ts
		return ts

	} else {
		return BetweenTimes(ts.mark, other.End())
	}
}

// RFC5545DateTimeLayout is the format string used by iCalendar (RFC5545). Note
// that "Z" is to be appended when the time is UTC.
//
// No dashes are used; this follows ISO-8601 Basic Format practice.
const RFC5545DateTimeLayout = "20060102T150405"

// RFC5545DateTimeZulu is the UTC format string used by iCalendar (RFC5545). Note
// that this cannot be used for parsing with time.Parse.
//
// No dashes are used; this follows ISO-8601 Basic Format practice.
const RFC5545DateTimeZulu = RFC5545DateTimeLayout + "Z"

func layoutHasTimezone(layout string) bool {
	return strings.IndexByte(layout, 'Z') >= 0 || strings.Contains(layout, "-07")
}

// Equal reports whether ts and us represent the same time start and duration.
// Two times can be equal even if they are in different locations.
// For example, 6:00 +0200 CEST and 4:00 UTC are Equal.
func (ts TimeSpan) Equal(us TimeSpan) bool {
	return ts.Duration() == us.Duration() && ts.Start().Equal(us.Start())
}
