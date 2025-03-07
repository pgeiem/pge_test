package engine

import (
	"fmt"
	"testing"
	"time"
)

func TestSegment_Duration(t *testing.T) {
	testCases := []struct {
		name     string
		start    time.Time
		end      time.Time
		expected time.Duration
	}{
		{
			name:     "OneHourDuration",
			start:    time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2023, 10, 1, 1, 0, 0, 0, time.UTC),
			expected: time.Hour,
		},
		{
			name:     "HalfHourDuration",
			start:    time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2023, 10, 1, 0, 30, 0, 0, time.UTC),
			expected: 30 * time.Minute,
		},
		{
			name:     "ZeroDuration",
			start:    time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			expected: 0,
		},
		{
			name:     "OneDayDuration",
			start:    time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC),
			expected: 24 * time.Hour,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			segment := Segment{Start: tc.start, End: tc.end}
			if segment.Duration() != tc.expected {
				t.Errorf("Duration() expected %v, got %v", tc.expected, segment.Duration())
			}
		})
	}
}

func TestSegment_IsWithin(t *testing.T) {
	testCases := []struct {
		name     string
		start    time.Time
		end      time.Time
		time     time.Time
		expected bool
	}{
		{
			name:     "WithinSegment",
			start:    time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2023, 10, 1, 1, 0, 0, 0, time.UTC),
			time:     time.Date(2023, 10, 1, 0, 30, 0, 0, time.UTC),
			expected: true,
		},
		{
			name:     "AfterSegment",
			start:    time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2023, 10, 1, 1, 0, 0, 0, time.UTC),
			time:     time.Date(2023, 10, 1, 1, 30, 0, 0, time.UTC),
			expected: false,
		},
		{
			name:     "AtSegmentStart",
			start:    time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2023, 10, 1, 1, 0, 0, 0, time.UTC),
			time:     time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			expected: true,
		},
		{
			name:     "AtSegmentEnd",
			start:    time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2023, 10, 1, 1, 0, 0, 0, time.UTC),
			time:     time.Date(2023, 10, 1, 1, 0, 0, 0, time.UTC),
			expected: false,
		},
		{
			name:     "BeforeSegment",
			start:    time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2023, 10, 1, 1, 0, 0, 0, time.UTC),
			time:     time.Date(2023, 9, 30, 23, 59, 59, 0, time.UTC),
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			segment := Segment{Start: tc.start, End: tc.end}
			if segment.IsWithin(tc.time) != tc.expected {
				t.Errorf("IsWithin(%v) expected %v, got %v", tc.time, tc.expected, segment.IsWithin(tc.time))
			}
		})
	}
}

func TestSegment_ToRelativeTimeSpan(t *testing.T) {
	testCases := []struct {
		name     string
		start    time.Time
		end      time.Time
		now      time.Time
		expected RelativeTimeSpan
	}{
		{
			name:     "FutureSegment",
			start:    time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2023, 10, 3, 0, 0, 0, 0, time.UTC),
			now:      time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			expected: RelativeTimeSpan{From: 24 * time.Hour, To: 48 * time.Hour},
		},
		{
			name:     "PastSegment",
			start:    time.Date(2023, 9, 30, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			now:      time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC),
			expected: RelativeTimeSpan{From: -48 * time.Hour, To: -24 * time.Hour},
		},
		{
			name:     "CurrentSegment",
			start:    time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC),
			now:      time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC),
			expected: RelativeTimeSpan{From: -12 * time.Hour, To: 12 * time.Hour},
		},
		{
			name:     "ZeroDurationSegment",
			start:    time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			now:      time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			expected: RelativeTimeSpan{From: 0, To: 0},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			segment := Segment{Start: tc.start, End: tc.end}
			relativeTimeSpan := segment.ToRelativeTimeSpan(tc.now)
			if relativeTimeSpan != tc.expected {
				t.Errorf("ToRelativeTimeSpan(%v) expected %v, got %v", tc.now, tc.expected, relativeTimeSpan)
			}
		})
	}
}

func TestRecurrentSegment_NextAndPrev(t *testing.T) {
	testCases := []struct {
		name         string
		startPattern string
		endPattern   string
		now          time.Time
		expectedNext Segment
		expectedPrev Segment
	}{
		{
			name:         "DailyPattern",
			startPattern: "pattern(2023/10/* 12:00:00)",
			endPattern:   "pattern(2023/10/* 18:00:00)",
			now:          time.Date(2023, 10, 15, 0, 0, 0, 0, time.UTC),
			expectedNext: Segment{Start: time.Date(2023, 10, 15, 12, 0, 0, 0, time.UTC), End: time.Date(2023, 10, 15, 18, 0, 0, 0, time.UTC)},
			expectedPrev: Segment{Start: time.Date(2023, 10, 14, 12, 0, 0, 0, time.UTC), End: time.Date(2023, 10, 14, 18, 0, 0, 0, time.UTC)},
		},
		{
			name:         "AtStartPattern",
			startPattern: "pattern(2023/10/* 12:00:00)",
			endPattern:   "pattern(2023/10/* 18:00:00)",
			now:          time.Date(2023, 10, 15, 12, 0, 0, 0, time.UTC),
			expectedNext: Segment{Start: time.Date(2023, 10, 16, 12, 0, 0, 0, time.UTC), End: time.Date(2023, 10, 16, 18, 0, 0, 0, time.UTC)},
			expectedPrev: Segment{Start: time.Date(2023, 10, 14, 12, 0, 0, 0, time.UTC), End: time.Date(2023, 10, 14, 18, 0, 0, 0, time.UTC)},
		},
		{
			name:         "AtEndPattern",
			startPattern: "pattern(2023/10/* 12:00:00)",
			endPattern:   "pattern(2023/10/* 18:00:00)",
			now:          time.Date(2023, 10, 15, 18, 0, 0, 0, time.UTC),
			expectedNext: Segment{Start: time.Date(2023, 10, 16, 12, 0, 0, 0, time.UTC), End: time.Date(2023, 10, 16, 18, 0, 0, 0, time.UTC)},
			expectedPrev: Segment{Start: time.Date(2023, 10, 15, 12, 0, 0, 0, time.UTC), End: time.Date(2023, 10, 15, 18, 0, 0, 0, time.UTC)},
		},
		{
			name:         "MonthlyPattern",
			startPattern: "pattern(2023/*/7 08:00:00)",
			endPattern:   "pattern(2023/*/9 10:00:00)",
			now:          time.Date(2023, 10, 3, 0, 0, 0, 0, time.UTC),
			expectedNext: Segment{Start: time.Date(2023, 10, 7, 8, 0, 0, 0, time.UTC), End: time.Date(2023, 10, 9, 10, 0, 0, 0, time.UTC)},
			expectedPrev: Segment{Start: time.Date(2023, 9, 7, 8, 0, 0, 0, time.UTC), End: time.Date(2023, 9, 9, 10, 0, 0, 0, time.UTC)},
		},
		{
			name:         "WeeklyPattern",
			startPattern: "pattern(2023/*/* WED 14:35:00)",
			endPattern:   "pattern(2023/*/* WED 16:50:00)",
			now:          time.Date(2023, 4, 1, 0, 0, 0, 0, time.UTC),
			expectedNext: Segment{Start: time.Date(2023, 4, 5, 14, 35, 0, 0, time.UTC), End: time.Date(2023, 4, 5, 16, 50, 0, 0, time.UTC)},
			expectedPrev: Segment{Start: time.Date(2023, 3, 29, 14, 35, 0, 0, time.UTC), End: time.Date(2023, 3, 29, 16, 50, 0, 0, time.UTC)},
		},
		{
			name:         "WeeklyPatternWithDuration",
			startPattern: "pattern(2023/*/* WED 14:35:00)",
			endPattern:   "duration(1h30m)",
			now:          time.Date(2023, 4, 1, 0, 0, 0, 0, time.UTC),
			expectedNext: Segment{Start: time.Date(2023, 4, 5, 14, 35, 0, 0, time.UTC), End: time.Date(2023, 4, 5, 16, 05, 0, 0, time.UTC)},
			expectedPrev: Segment{Start: time.Date(2023, 3, 29, 14, 35, 0, 0, time.UTC), End: time.Date(2023, 3, 29, 16, 05, 0, 0, time.UTC)},
		},
		{
			name:         "Periodic8Hours",
			startPattern: "periodic(8h)",
			endPattern:   "duration(1h30m)",
			now:          time.Date(2023, 4, 1, 0, 0, 0, 0, time.UTC),
			expectedNext: Segment{Start: time.Date(2023, 4, 1, 8, 0, 0, 0, time.UTC), End: time.Date(2023, 4, 1, 9, 30, 0, 0, time.UTC)},
			expectedPrev: Segment{Start: time.Date(2023, 3, 31, 16, 0, 0, 0, time.UTC), End: time.Date(2023, 3, 31, 17, 30, 0, 0, time.UTC)},
		},
		{
			name:         "Periodic24HoursWithPattern",
			startPattern: "periodic(24h)",
			endPattern:   "pattern(2023/*/7 08:00:00)",
			now:          time.Date(2023, 4, 1, 0, 0, 0, 0, time.UTC),
			expectedNext: Segment{Start: time.Date(2023, 4, 2, 0, 0, 0, 0, time.UTC), End: time.Date(2023, 4, 7, 8, 0, 0, 0, time.UTC)},
			expectedPrev: Segment{Start: time.Date(2023, 3, 31, 0, 0, 0, 0, time.UTC), End: time.Date(2023, 4, 7, 8, 0, 0, 0, time.UTC)},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Printf("\n------\nHandling test case %v\n", tc)

			start, errStart := ParseRecurrentDate(tc.startPattern)
			if errStart != nil {
				t.Fatalf("unexpected error in ParseRecurrentDate for start: %v", errStart)
			}
			end, errEnd := ParseRecurrentDate(tc.endPattern)
			if errEnd != nil {
				t.Fatalf("unexpected error in ParseRecurrentDate for end: %v", errEnd)
			}
			rs := RecurrentSegment{Start: start, End: end}

			// Test Next
			segmentNext, errNext := rs.Next(tc.now)
			if errNext != nil {
				t.Errorf("unexpected error in Next: %v", errNext)
			}
			if !segmentNext.Start.Equal(tc.expectedNext.Start) || !segmentNext.End.Equal(tc.expectedNext.End) {
				t.Errorf("next(%v) expected segment %v in Next, got %v", tc.now, tc.expectedNext, segmentNext)
			}

			// Test Prev
			segmentPrev, errPrev := rs.Prev(tc.now)
			if errPrev != nil {
				t.Errorf("unexpected error in Prev: %v", errPrev)
			}
			if !segmentPrev.Start.Equal(tc.expectedPrev.Start) || !segmentPrev.End.Equal(tc.expectedPrev.End) {
				t.Errorf("Prev(%v) expected segment %v in Prev, got %v", tc.now, tc.expectedPrev, segmentPrev)
			}
		})
	}
}

func TestRecurrentSegment_Between(t *testing.T) {
	testCases := []struct {
		name         string
		startPattern string
		endPattern   string
		from         time.Time
		to           time.Time
		expected     []Segment
	}{
		{
			name:         "DailyPatternBetween",
			startPattern: "pattern(2023/10/* 12:00:00)",
			endPattern:   "pattern(2023/10/* 18:00:00)",
			from:         time.Date(2023, 10, 14, 0, 0, 0, 0, time.UTC),
			to:           time.Date(2023, 10, 16, 0, 0, 0, 0, time.UTC),
			expected: []Segment{
				{Start: time.Date(2023, 10, 14, 12, 0, 0, 0, time.UTC), End: time.Date(2023, 10, 14, 18, 0, 0, 0, time.UTC)},
				{Start: time.Date(2023, 10, 15, 12, 0, 0, 0, time.UTC), End: time.Date(2023, 10, 15, 18, 0, 0, 0, time.UTC)},
			},
		},
		{
			name:         "MonthlyPatternBetween",
			startPattern: "pattern(2023/*/7 08:00:00)",
			endPattern:   "pattern(2023/*/9 10:00:00)",
			from:         time.Date(2023, 9, 6, 0, 0, 0, 0, time.UTC),
			to:           time.Date(2023, 12, 10, 0, 0, 0, 0, time.UTC),
			expected: []Segment{
				{Start: time.Date(2023, 9, 7, 8, 0, 0, 0, time.UTC), End: time.Date(2023, 9, 9, 10, 0, 0, 0, time.UTC)},
				{Start: time.Date(2023, 10, 7, 8, 0, 0, 0, time.UTC), End: time.Date(2023, 10, 9, 10, 0, 0, 0, time.UTC)},
				{Start: time.Date(2023, 11, 7, 8, 0, 0, 0, time.UTC), End: time.Date(2023, 11, 9, 10, 0, 0, 0, time.UTC)},
				{Start: time.Date(2023, 12, 7, 8, 0, 0, 0, time.UTC), End: time.Date(2023, 12, 9, 10, 0, 0, 0, time.UTC)},
			},
		},
		{
			name:         "WeeklyPatternBetween",
			startPattern: "pattern(2023/*/* WED 14:35:00)",
			endPattern:   "pattern(2023/*/* WED 16:50:00)",
			from:         time.Date(2023, 3, 28, 0, 0, 0, 0, time.UTC),
			to:           time.Date(2023, 4, 12, 23, 0, 0, 0, time.UTC),
			expected: []Segment{
				{Start: time.Date(2023, 3, 29, 14, 35, 0, 0, time.UTC), End: time.Date(2023, 3, 29, 16, 50, 0, 0, time.UTC)},
				{Start: time.Date(2023, 4, 5, 14, 35, 0, 0, time.UTC), End: time.Date(2023, 4, 5, 16, 50, 0, 0, time.UTC)},
				{Start: time.Date(2023, 4, 12, 14, 35, 0, 0, time.UTC), End: time.Date(2023, 4, 12, 16, 50, 0, 0, time.UTC)},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Printf("\n------\nHandling test case %v\n", tc)

			start, errStart := ParseRecurrentDate(tc.startPattern)
			if errStart != nil {
				t.Fatalf("unexpected error in ParseRecurrentDate for start: %v", errStart)
			}
			end, errEnd := ParseRecurrentDate(tc.endPattern)
			if errEnd != nil {
				t.Fatalf("unexpected error in ParseRecurrentDate for end: %v", errEnd)
			}
			rs := RecurrentSegment{Start: start, End: end}

			segments := rs.Between(tc.from, tc.to)

			if len(segments) != len(tc.expected) {
				t.Errorf("Between(%v, %v) expected %v segments, got %v", tc.from, tc.to, len(tc.expected), len(segments))
			} else {
				for i, segment := range segments {
					if !segment.Start.Equal(tc.expected[i].Start) || !segment.End.Equal(tc.expected[i].End) {
						t.Errorf("Between(%v, %v) expected segment %v, got %v", tc.from, tc.to, tc.expected[i], segment)
					}
				}
			}
		})
	}
}

func TestRecurrentSegment_IsWithin(t *testing.T) {
	testCases := []struct {
		name         string
		startPattern string
		endPattern   string
		time         time.Time
		expected     bool
		expectedSeg  Segment
	}{
		{
			name:         "WithinRecurrentSegment",
			startPattern: "pattern(2023/10/* 12:00:00)",
			endPattern:   "pattern(2023/10/* 18:00:00)",
			time:         time.Date(2023, 10, 15, 14, 0, 0, 0, time.UTC),
			expected:     true,
			expectedSeg:  Segment{Start: time.Date(2023, 10, 15, 12, 0, 0, 0, time.UTC), End: time.Date(2023, 10, 15, 18, 0, 0, 0, time.UTC)},
		},
		{
			name:         "BeforeRecurrentSegment",
			startPattern: "pattern(2023/10/* 12:00:00)",
			endPattern:   "pattern(2023/10/* 18:00:00)",
			time:         time.Date(2023, 10, 15, 10, 0, 0, 0, time.UTC),
			expected:     false,
			expectedSeg:  Segment{},
		},
		{
			name:         "AfterRecurrentSegment",
			startPattern: "pattern(2023/10/* 12:00:00)",
			endPattern:   "pattern(2023/10/* 18:00:00)",
			time:         time.Date(2023, 10, 15, 19, 0, 0, 0, time.UTC),
			expected:     false,
			expectedSeg:  Segment{},
		},
		{
			name:         "AtRecurrentSegmentStart",
			startPattern: "pattern(2023/10/* 12:00:00)",
			endPattern:   "pattern(2023/10/* 18:00:00)",
			time:         time.Date(2023, 10, 15, 12, 0, 0, 0, time.UTC),
			expected:     true,
			expectedSeg:  Segment{Start: time.Date(2023, 10, 15, 12, 0, 0, 0, time.UTC), End: time.Date(2023, 10, 15, 18, 0, 0, 0, time.UTC)},
		},
		{
			name:         "AtRecurrentSegmentEnd",
			startPattern: "pattern(2023/10/* 12:00:00)",
			endPattern:   "pattern(2023/10/* 18:00:00)",
			time:         time.Date(2023, 10, 15, 18, 0, 0, 0, time.UTC),
			expected:     false,
			expectedSeg:  Segment{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			start, errStart := ParseRecurrentDate(tc.startPattern)
			if errStart != nil {
				t.Fatalf("unexpected error in ParseRecurrentDate for start: %v", errStart)
			}
			end, errEnd := ParseRecurrentDate(tc.endPattern)
			if errEnd != nil {
				t.Fatalf("unexpected error in ParseRecurrentDate for end: %v", errEnd)
			}
			rs := RecurrentSegment{Start: start, End: end}

			isWithin, segment, err := rs.IsWithinWithSegment(tc.time)
			if err != nil {
				t.Errorf("unexpected error in IsWithinWithSegment: %v", err)
			}
			if isWithin != tc.expected {
				t.Errorf("IsWithinWithSegment(%v) expected %v, got %v", tc.time, tc.expected, isWithin)
			}
			if !segment.Start.Equal(tc.expectedSeg.Start) || !segment.End.Equal(tc.expectedSeg.End) {
				t.Errorf("IsWithinWithSegment(%v) expected segment %v, got %v", tc.time, tc.expectedSeg, segment)
			}
		})
	}
}
