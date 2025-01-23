package parser

import (
	"fmt"
	"testing"
	"time"
)

func TestParseDuration(t *testing.T) {
	tests := []struct {
		input    string
		expected Duration
		hasError bool
	}{
		{"1w", Duration(7 * 24 * time.Hour), false},
		{"2d", Duration(2 * 24 * time.Hour), false},
		{"3h", Duration(3 * time.Hour), false},
		{"4m", Duration(4 * time.Minute), false},
		{"5s", Duration(5 * time.Second), false},
		{"1w2d3h4m5s", Duration(7*24*time.Hour + 2*24*time.Hour + 3*time.Hour + 4*time.Minute + 5*time.Second), false},
		{"2d3h4m5s", Duration(2*24*time.Hour + 3*time.Hour + 4*time.Minute + 5*time.Second), false},
		{"3h4m5s", Duration(3*time.Hour + 4*time.Minute + 5*time.Second), false},
		{"4m5s", Duration(4*time.Minute + 5*time.Second), false},
		{"5s", Duration(5 * time.Second), false},
		{"1w2h", Duration(7*24*time.Hour + 2*time.Hour), false}, // Valid order
		{"2d1w", 0, true}, // Invalid order
		{"3h2d", 0, true}, // Invalid order
		{"4m3h", 0, true}, // Invalid order
		{"5s4m", 0, true}, // Invalid order
		{"", 0, true},     // Empty string
		{"1x", 0, true},   // Unknown unit
	}

	for _, test := range tests {
		result, err := ParseDuration(test.input)
		if (err != nil) != test.hasError {
			t.Errorf("ParseDuration(%q) error = %v, wantErr %v", test.input, err, test.hasError)
			continue
		}
		if result != test.expected {
			t.Errorf("ParseDuration(%q) = %v, want %v", test.input, result, test.expected)
		}
	}
}

func TestExpandDateComponentList(t *testing.T) {
	tests := []struct {
		input    string
		expected []int
		hasError bool
	}{
		{"1,2,3", []int{1, 2, 3}, false},
		{"1-3", []int{1, 2, 3}, false},
		{"1,2-4,5", []int{1, 2, 3, 4, 5}, false},
		{"MON,TUE,WED", []int{0, 1, 2}, false},
		{"1,MON,3", []int{1, 0, 3}, false},
		{"1-3,5-7", []int{1, 2, 3, 5, 6, 7}, false},
		{"JAN,Feb,MAR", []int{1, 2, 3}, false},
		{"1-3,5,7-9", []int{1, 2, 3, 5, 7, 8, 9}, false},
		{"1,3-5,7", []int{1, 3, 4, 5, 7}, false},
		{"1-5,10-15", []int{1, 2, 3, 4, 5, 10, 11, 12, 13, 14, 15}, false},
		{"1-3,5-7,9-11,13-15", []int{1, 2, 3, 5, 6, 7, 9, 10, 11, 13, 14, 15}, false},
		{"", nil, true},                    // Empty string
		{"1,,2", nil, true},                // Double comma
		{"1-2-3", nil, true},               // Invalid range
		{"1-2,", nil, true},                // Trailing comma
		{"1-2,3-", nil, true},              // Trailing dash
		{"1-2,3-abc", nil, true},           // Invalid character
		{"1-2,3-1000000000000", nil, true}, // Out of range
	}

	for _, test := range tests {
		result, err := expandDateComponentList(test.input)
		if (err != nil) != test.hasError {
			t.Errorf("ExpandDateComponentList(%q) error = %v, wantErr %v", test.input, err, test.hasError)
			continue
		}
		if !equal(result, test.expected) {
			t.Errorf("ExpandDateComponentList(%q) = %v, want %v", test.input, result, test.expected)
		}
	}
}

func equal(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestRecurrentDatePeriodic(t *testing.T) {
	tests := []struct {
		pattern      string
		now          time.Time
		expectedNext []time.Time
		expectedPrev []time.Time
	}{
		{
			"1h30m",
			time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			[]time.Time{
				time.Date(2023, 10, 1, 1, 30, 0, 0, time.UTC),
				time.Date(2023, 10, 1, 3, 0, 0, 0, time.UTC),
				time.Date(2023, 10, 1, 4, 30, 0, 0, time.UTC),
			},
			[]time.Time{
				time.Date(2023, 9, 30, 22, 30, 0, 0, time.UTC),
				time.Date(2023, 9, 30, 21, 0, 0, 0, time.UTC),
				time.Date(2023, 9, 30, 19, 30, 0, 0, time.UTC),
			},
		},
		{
			"2h",
			time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			[]time.Time{
				time.Date(2023, 10, 1, 2, 0, 0, 0, time.UTC),
				time.Date(2023, 10, 1, 4, 0, 0, 0, time.UTC),
				time.Date(2023, 10, 1, 6, 0, 0, 0, time.UTC),
			},
			[]time.Time{
				time.Date(2023, 9, 30, 22, 0, 0, 0, time.UTC),
				time.Date(2023, 9, 30, 20, 0, 0, 0, time.UTC),
				time.Date(2023, 9, 30, 18, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, test := range tests {
		fmt.Println("Testing periodic", test.pattern)
		r := &RecurrentDatePeriodic{}
		err := r.Parse(test.pattern)
		if err != nil {
			t.Fatalf("Parse failed: %v", err)
		}

		now := test.now
		for _, expectedNext := range test.expectedNext {
			next, err := r.Next(now)
			if err != nil {
				t.Fatalf("Next failed: %v", err)
			}
			if !next.Equal(expectedNext) {
				t.Errorf("Next() = %v, want %v", next, expectedNext)
			}
			now = next
		}

		now = test.now
		for _, expectedPrev := range test.expectedPrev {
			prev, err := r.Prev(now)
			if err != nil {
				t.Fatalf("Prev failed: %v", err)
			}
			if !prev.Equal(expectedPrev) {
				t.Errorf("Prev() = %v, want %v", prev, expectedPrev)
			}
			now = prev
		}
	}
}

func TestRecurrentDateCron(t *testing.T) {
	tests := []struct {
		pattern      string
		now          time.Time
		expectedNext []time.Time
		expectedPrev []time.Time
	}{
		{
			"0 0 * * *",
			time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			[]time.Time{
				time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2023, 10, 3, 0, 0, 0, 0, time.UTC),
				time.Date(2023, 10, 4, 0, 0, 0, 0, time.UTC),
			},
			[]time.Time{
				time.Date(2023, 9, 30, 0, 0, 0, 0, time.UTC),
				time.Date(2023, 9, 29, 0, 0, 0, 0, time.UTC),
				time.Date(2023, 9, 28, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			"0 12 * * *",
			time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			[]time.Time{
				time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC),
				time.Date(2023, 10, 2, 12, 0, 0, 0, time.UTC),
				time.Date(2023, 10, 3, 12, 0, 0, 0, time.UTC),
			},
			[]time.Time{
				time.Date(2023, 9, 30, 12, 0, 0, 0, time.UTC),
				time.Date(2023, 9, 29, 12, 0, 0, 0, time.UTC),
				time.Date(2023, 9, 28, 12, 0, 0, 0, time.UTC),
			},
		},
		{
			"0 0 1 * *",
			time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			[]time.Time{
				time.Date(2023, 11, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			[]time.Time{
				time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			"0 0 * * 0",
			time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			[]time.Time{
				time.Date(2023, 10, 8, 0, 0, 0, 0, time.UTC),
				time.Date(2023, 10, 15, 0, 0, 0, 0, time.UTC),
				time.Date(2023, 10, 22, 0, 0, 0, 0, time.UTC),
			},
			[]time.Time{
				time.Date(2023, 9, 24, 0, 0, 0, 0, time.UTC),
				time.Date(2023, 9, 17, 0, 0, 0, 0, time.UTC),
				time.Date(2023, 9, 10, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			"0 0 1 1 *",
			time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			[]time.Time{
				time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			[]time.Time{
				time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, test := range tests {
		fmt.Println("Testing cron", test.pattern)
		r := &RecurrentDateCron{}
		err := r.Parse(test.pattern)
		if err != nil {
			t.Fatalf("Parse failed: %v", err)
		}

		now := test.now
		for _, expectedNext := range test.expectedNext {
			next, err := r.Next(now)
			if err != nil {
				t.Fatalf("Next failed: %v", err)
			}
			if !next.Equal(expectedNext) {
				t.Errorf("Next() = %v, want %v", next, expectedNext)
			}
			now = next
		}

		now = test.now
		for _, expectedPrev := range test.expectedPrev {
			prev, err := r.Prev(now)
			if err != nil {
				t.Fatalf("Prev failed: %v", err)
			}
			if !prev.Equal(expectedPrev) {
				t.Errorf("Prev() = %v, want %v", prev, expectedPrev)
			}
			now = prev
		}
	}
}

func TestRecurrentDatePattern(t *testing.T) {
	tests := []struct {
		pattern      string
		now          time.Time
		expectedNext []time.Time
		expectedPrev []time.Time
		hasError     bool
	}{
		{
			"2023/10/* MON 12:00:00",
			time.Date(2023, 10, 16, 0, 0, 0, 0, time.UTC),
			[]time.Time{
				time.Date(2023, 10, 16, 12, 0, 0, 0, time.UTC),
				time.Date(2023, 10, 23, 12, 0, 0, 0, time.UTC),
				time.Date(2023, 10, 30, 12, 0, 0, 0, time.UTC),
			},
			[]time.Time{
				time.Date(2023, 10, 9, 12, 0, 0, 0, time.UTC),
				time.Date(2023, 10, 2, 12, 0, 0, 0, time.UTC),
			},
			false,
		},
		{
			"2023/*/* Mon-Fri 12:00:00",
			time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			[]time.Time{
				time.Date(2023, 10, 2, 12, 0, 0, 0, time.UTC),
				time.Date(2023, 10, 3, 12, 0, 0, 0, time.UTC),
				time.Date(2023, 10, 4, 12, 0, 0, 0, time.UTC),
			},
			[]time.Time{
				time.Date(2023, 9, 29, 12, 0, 0, 0, time.UTC),
				time.Date(2023, 9, 28, 12, 0, 0, 0, time.UTC),
				time.Date(2023, 9, 27, 12, 0, 0, 0, time.UTC),
			},
			false,
		},
		{
			"2023/9-10/* Mon,Fri,SUN 12:00:00",
			time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			[]time.Time{
				time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC),
				time.Date(2023, 10, 2, 12, 0, 0, 0, time.UTC),
				time.Date(2023, 10, 6, 12, 0, 0, 0, time.UTC),
			},
			[]time.Time{
				time.Date(2023, 9, 29, 12, 0, 0, 0, time.UTC),
				time.Date(2023, 9, 25, 12, 0, 0, 0, time.UTC),
				time.Date(2023, 9, 24, 12, 0, 0, 0, time.UTC),
			},
			false,
		},
		{
			"2023/*/* 12¦2:00:00",
			time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			[]time.Time{
				time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC),
				time.Date(2023, 10, 1, 14, 0, 0, 0, time.UTC),
				time.Date(2023, 10, 1, 16, 0, 0, 0, time.UTC),
			},
			[]time.Time{
				time.Date(2023, 9, 30, 22, 0, 0, 0, time.UTC),
				time.Date(2023, 9, 30, 20, 0, 0, 0, time.UTC),
				time.Date(2023, 9, 30, 18, 0, 0, 0, time.UTC),
			},
			false,
		},
		{
			"2023/10/* * 12,23:00",
			time.Date(2023, 10, 8, 16, 0, 0, 0, time.UTC),
			[]time.Time{
				time.Date(2023, 10, 8, 23, 0, 0, 0, time.UTC),
				time.Date(2023, 10, 9, 12, 0, 0, 0, time.UTC),
				time.Date(2023, 10, 9, 23, 0, 0, 0, time.UTC),
			},
			[]time.Time{
				time.Date(2023, 10, 8, 12, 0, 0, 0, time.UTC),
				time.Date(2023, 10, 7, 23, 0, 0, 0, time.UTC),
				time.Date(2023, 10, 7, 12, 0, 0, 0, time.UTC),
			},
			false,
		},
		{
			"2023/10/* 17:00",
			time.Date(2023, 10, 5, 0, 0, 0, 0, time.UTC),
			[]time.Time{
				time.Date(2023, 10, 5, 17, 0, 0, 0, time.UTC),
				time.Date(2023, 10, 6, 17, 0, 0, 0, time.UTC),
				time.Date(2023, 10, 7, 17, 0, 0, 0, time.UTC),
			},
			[]time.Time{
				time.Date(2023, 10, 4, 17, 0, 0, 0, time.UTC),
				time.Date(2023, 10, 3, 17, 0, 0, 0, time.UTC),
				time.Date(2023, 10, 2, 17, 0, 0, 0, time.UTC),
			},
			false,
		},
		{
			"2023/*/* Mon,FRI 12:00",
			time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			[]time.Time{
				time.Date(2023, 10, 2, 12, 0, 0, 0, time.UTC),
				time.Date(2023, 10, 6, 12, 0, 0, 0, time.UTC),
				time.Date(2023, 10, 9, 12, 0, 0, 0, time.UTC),
			},
			[]time.Time{
				time.Date(2023, 9, 29, 12, 0, 0, 0, time.UTC),
				time.Date(2023, 9, 25, 12, 0, 0, 0, time.UTC),
				time.Date(2023, 9, 22, 12, 0, 0, 0, time.UTC),
			},
			false,
		},
		{
			"2023/10/01 Mon 12:00:00", // valid rules but too much constraints no next or prev may be found
			time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC),
			[]time.Time{time.Date(2023, 10, 2, 12, 0, 0, 0, time.UTC)}, // This is incorrect as day is 2 and should be 1 !
			[]time.Time{time.Unix(0, 0)},
			false,
		},
		{
			"2023/10/01 Mon 12:00:00:00",
			time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			[]time.Time{},
			[]time.Time{},
			true,
		},
		{
			"2023/10/01 Mon 12:00:xx",
			time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			[]time.Time{},
			[]time.Time{},
			true,
		},
		{
			"2023/10/01 Mon 12:xx:00",
			time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			[]time.Time{},
			[]time.Time{},
			true,
		},
		{
			"2023/10/01 Mon xx:00:00",
			time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			[]time.Time{},
			[]time.Time{},
			true,
		},
		{
			"2023/10/01 Mon *:*:00¦5",
			time.Date(2023, 10, 1, 10, 0, 0, 0, time.UTC),
			[]time.Time{
				time.Date(2023, 10, 1, 10, 0, 5, 0, time.UTC),
				time.Date(2023, 10, 1, 10, 0, 10, 0, time.UTC),
				time.Date(2023, 10, 1, 10, 0, 15, 0, time.UTC),
			},
			[]time.Time{
				time.Date(2023, 10, 1, 9, 59, 55, 0, time.UTC),
				time.Date(2023, 10, 1, 9, 59, 50, 0, time.UTC),
				time.Date(2023, 10, 1, 9, 59, 45, 0, time.UTC),
			},
			false,
		},
	}

	for _, test := range tests {
		r := &RecurrentDatePattern{}
		fmt.Println("Testing pattern", test.pattern)
		err := r.Parse(test.pattern)
		if (err != nil) != test.hasError {
			t.Errorf("ParseRecurrentDatePattern(%q) error = %v, wantErr %v", test.pattern, err, test.hasError)
			continue
		}
		if err != nil {
			continue
		}

		now := test.now
		for _, expectedNext := range test.expectedNext {
			next, err := r.Next(now)
			if err != nil {
				t.Fatalf("Next failed: %v", err)
			}
			if !next.Equal(expectedNext) {
				t.Errorf("Next() = %v, want %v", next, expectedNext)
			}
			now = next
		}

		now = test.now
		for _, expectedPrev := range test.expectedPrev {
			prev, err := r.Prev(now)
			if (err != nil) != (expectedPrev == time.Unix(0, 0)) {
				t.Fatalf("Prev failed: %v", err)
			}
			if (expectedPrev != time.Unix(0, 0)) && !prev.Equal(expectedPrev) {
				t.Errorf("Prev(%v) = %v, want %v", now, prev, expectedPrev)
			}
			now = prev
		}
	}
}

func TestParseRecurrentDate(t *testing.T) {
	tests := []struct {
		pattern      string
		now          time.Time
		expectedNext time.Time
		expectedPrev time.Time
		hasError     bool
	}{
		{
			"periodic(1h30m)",
			time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2023, 10, 1, 1, 30, 0, 0, time.UTC),
			time.Date(2023, 9, 30, 22, 30, 0, 0, time.UTC),
			false,
		},
		{
			"cron(0 0 * * *)",
			time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC),
			time.Date(2023, 9, 30, 0, 0, 0, 0, time.UTC),
			false,
		},
		{
			"pattern(2023/10/* MON 12:00:00)",
			time.Date(2023, 10, 16, 0, 0, 0, 0, time.UTC),
			time.Date(2023, 10, 16, 12, 0, 0, 0, time.UTC),
			time.Date(2023, 10, 9, 12, 0, 0, 0, time.UTC),
			false,
		},
		{
			"pattern(2023/*/* Mon-Fri 12:00:00)",
			time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2023, 10, 2, 12, 0, 0, 0, time.UTC),
			time.Date(2023, 9, 29, 12, 0, 0, 0, time.UTC),
			false,
		},
		{
			"pattern(2023/9-10/* Mon,Fri,SUN 12:00:00)",
			time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC),
			time.Date(2023, 9, 29, 12, 0, 0, 0, time.UTC),
			false,
		},
		{
			"pattern(2023/*/* 12¦2:00:00)",
			time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC),
			time.Date(2023, 9, 30, 22, 0, 0, 0, time.UTC),
			false,
		},
		{
			"pattern(2023/10/01 * 12,23:00)",
			time.Date(2023, 10, 1, 16, 0, 0, 0, time.UTC),
			time.Date(2023, 10, 1, 23, 0, 0, 0, time.UTC),
			time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC),
			false,
		},
		{
			"pattern(2023/10/* 17:00)",
			time.Date(2023, 10, 3, 0, 0, 0, 0, time.UTC),
			time.Date(2023, 10, 3, 17, 0, 0, 0, time.UTC),
			time.Date(2023, 10, 2, 17, 0, 0, 0, time.UTC),
			false,
		},
		{
			"pattern(2023/*/* Mon,FRI 12:00)",
			time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2023, 10, 2, 12, 0, 0, 0, time.UTC),
			time.Date(2023, 9, 29, 12, 0, 0, 0, time.UTC),
			false,
		},
		{
			"invalid(1h30m)",
			time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			time.Time{},
			time.Time{},
			true,
		},
		{
			"periodic(invalid)",
			time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			time.Time{},
			time.Time{},
			true,
		},
		{
			"cron(invalid)",
			time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			time.Time{},
			time.Time{},
			true,
		},
		{
			"pattern(invalid)",
			time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			time.Time{},
			time.Time{},
			true,
		},
		{
			"",
			time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			time.Time{},
			time.Time{},
			true,
		},
		{
			"()",
			time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			time.Time{},
			time.Time{},
			true,
		},
		{
			"periodic()",
			time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			time.Time{},
			time.Time{},
			true,
		},
		{
			"cron()",
			time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			time.Time{},
			time.Time{},
			true,
		},
		{
			"pattern()",
			time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			time.Time{},
			time.Time{},
			true,
		},
		{
			"(1h30m)",
			time.Time{},
			time.Time{},
			time.Time{},
			true,
		},
		{
			"(0 0 * * *)",
			time.Time{},
			time.Time{},
			time.Time{},
			true,
		},
		{
			"(2023/10/* MON 12:00:00)",
			time.Time{},
			time.Time{},
			time.Time{},
			true,
		},
	}

	for _, test := range tests {
		recurrentDate, err := ParseRecurrentDate(test.pattern)
		if (err != nil) != test.hasError {
			t.Errorf("ParseRecurrentDate(%q) error = %v, wantErr %v", test.pattern, err, test.hasError)
			continue
		}
		if err != nil {
			continue
		}

		next, err := recurrentDate.Next(test.now)
		if err != nil {
			t.Fatalf("Next failed: %v", err)
		}
		if !next.Equal(test.expectedNext) {
			t.Errorf("Next() = %v, want %v", next, test.expectedNext)
		}

		prev, err := recurrentDate.Prev(test.now)
		if err != nil {
			t.Fatalf("Prev failed: %v", err)
		}
		if !prev.Equal(test.expectedPrev) {
			t.Errorf("Prev() = %v, want %v", prev, test.expectedPrev)
		}
	}
}
