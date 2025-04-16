package engine

import (
	"regexp"
	"testing"
	"time"
)

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9-_()]+`)

func clearString(str string) string {
	return nonAlphanumericRegex.ReplaceAllString(str, "")
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
		{"", []int{}, false},               // Empty string
		{"1,,2", nil, true},                // Double comma
		{"1-2-3", nil, true},               // Invalid range
		{"1-2,", nil, true},                // Trailing comma
		{"1-2,3-", nil, true},              // Trailing dash
		{"1-2,3-abc", nil, true},           // Invalid character
		{"1-2,3-1000000000000", nil, true}, // Out of range
	}

	for _, test := range tests {
		t.Run(clearString(test.input), func(t *testing.T) {
			result, err := expandDateComponentList(test.input)
			if (err != nil) != test.hasError {
				t.Errorf("ExpandDateComponentList(%q) error = %v, wantErr %v", test.input, err, test.hasError)
				return
			}
			if !equal(result, test.expected) {
				t.Errorf("ExpandDateComponentList(%q) = %v, want %v", test.input, result, test.expected)
			}
		})
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
		pattern       string
		now           time.Time
		expectedFirst time.Time
		expectedNext  []time.Time
		expectedPrev  []time.Time
	}{
		{
			"1h30m",
			time.Date(2023, 10, 1, 0, 0, 0, 0, time.Local),
			time.Date(2023, 10, 1, 0, 0, 0, 0, time.Local),
			[]time.Time{
				time.Date(2023, 10, 1, 1, 30, 0, 0, time.Local),
				time.Date(2023, 10, 1, 3, 0, 0, 0, time.Local),
				time.Date(2023, 10, 1, 4, 30, 0, 0, time.Local),
			},
			[]time.Time{
				time.Date(2023, 9, 30, 22, 30, 0, 0, time.Local),
				time.Date(2023, 9, 30, 21, 0, 0, 0, time.Local),
				time.Date(2023, 9, 30, 19, 30, 0, 0, time.Local),
			},
		},
		{
			"2h",
			time.Date(2023, 10, 1, 0, 0, 0, 0, time.Local),
			time.Date(2023, 10, 1, 0, 0, 0, 0, time.Local),
			[]time.Time{
				time.Date(2023, 10, 1, 2, 0, 0, 0, time.Local),
				time.Date(2023, 10, 1, 4, 0, 0, 0, time.Local),
				time.Date(2023, 10, 1, 6, 0, 0, 0, time.Local),
			},
			[]time.Time{
				time.Date(2023, 9, 30, 22, 0, 0, 0, time.Local),
				time.Date(2023, 9, 30, 20, 0, 0, 0, time.Local),
				time.Date(2023, 9, 30, 18, 0, 0, 0, time.Local),
			},
		},
	}

	for _, test := range tests {
		t.Run(clearString(test.pattern), func(t *testing.T) {
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

			first, err := r.First(test.now)
			if err != nil {
				t.Fatalf("Next failed: %v", err)
			}
			if !first.Equal(test.expectedFirst) {
				t.Errorf("Next(%v) = %v, want %v", now, first, test.expectedFirst)
			}

		})
	}
}

func TestBuilRRuleFromDatePattern(t *testing.T) {

	tests := []struct {
		pattern       string
		expectedrrule string
		hasError      bool
	}{
		{
			"2023/10/* MON 12:00:00",
			"FREQ=DAILY;BYMONTH=10;BYDAY=MO;BYHOUR=12;BYMINUTE=0;BYSECOND=0",
			false,
		},
		{
			"2023/*/* Mon-Fri 12:00:00",
			"FREQ=DAILY;BYDAY=MO,TU,WE,TH,FR;BYHOUR=12;BYMINUTE=0;BYSECOND=0",
			false,
		},
		{
			"2023/9-10/* Mon,Fri,SUN 12:00:00",
			"FREQ=DAILY;BYMONTH=9,10;BYDAY=MO,FR,SU;BYHOUR=12;BYMINUTE=0;BYSECOND=0",
			false,
		},
		{
			"2023/*/* 12:00:00",
			"FREQ=DAILY;BYHOUR=12;BYMINUTE=0;BYSECOND=0",
			false,
		},
		{
			"2023/10/* * 12,23:00",
			"FREQ=DAILY;BYMONTH=10;BYHOUR=12,23;BYMINUTE=0",
			false,
		},
		{
			"2023/10/* 17:00",
			"FREQ=DAILY;BYMONTH=10;BYHOUR=17;BYMINUTE=0",
			false,
		},
		{
			"2023/*/* Mon,FRI 12:00",
			"FREQ=DAILY;BYDAY=MO,FR;BYHOUR=12;BYMINUTE=0",
			false,
		},
		{
			"2023/*/* Mon,FRI 12:00 COUNT=5",
			"FREQ=DAILY;COUNT=5;BYDAY=MO,FR;BYHOUR=12;BYMINUTE=0",
			false,
		},
		{
			"2023/10/01 Mon 12:00:00",
			"",
			true,
		},
		{
			"2023/10/01 Mon *:*:00",
			"FREQ=MINUTELY;BYMONTH=10;BYMONTHDAY=1;BYDAY=MO;BYSECOND=0",
			false,
		},
		{
			"2023/10/01 Mon 12:00:00:00",
			"",
			true,
		},
		{
			"2023/10/01 Mon 12:00:xx",
			"",
			true,
		},
		{
			"2023/10/01 Mon 12:xx:00",
			"",
			true,
		},
		{
			"2023/10/01 Mon xx:00:00",
			"",
			true,
		},
	}

	for _, test := range tests {
		t.Run(clearString(test.pattern), func(t *testing.T) {
			rrule, err := buildRRuleFromDatePattern(test.pattern)
			if (err != nil) != test.hasError {
				t.Errorf("BuilRRuleFromDatePattern(%q) error = %v, wantErr %v", test.pattern, err, test.hasError)
				return
			}
			if err != nil {
				return
			}
			if rrule.String() != test.expectedrrule {
				t.Errorf("BuilRRuleFromDatePattern(%q) = %v, want %v", test.pattern, rrule, test.expectedrrule)
			}
		})
	}
}

func TestRecurrentDatePattern(t *testing.T) {
	tests := []struct {
		pattern       string
		now           time.Time
		expectedFirst time.Time
		expectedNext  []time.Time
		expectedPrev  []time.Time
		hasError      bool
	}{
		{
			"2023/10/* MON 12:00:00",
			time.Date(2023, 10, 16, 0, 0, 0, 0, time.Local),
			time.Date(2023, 10, 16, 12, 0, 0, 0, time.Local),
			[]time.Time{
				time.Date(2023, 10, 16, 12, 0, 0, 0, time.Local),
				time.Date(2023, 10, 23, 12, 0, 0, 0, time.Local),
				time.Date(2023, 10, 30, 12, 0, 0, 0, time.Local),
			},
			[]time.Time{
				time.Date(2023, 10, 9, 12, 0, 0, 0, time.Local),
				time.Date(2023, 10, 2, 12, 0, 0, 0, time.Local),
			},
			false,
		},
		{
			"2023/*/* Mon-Fri 12:00:00",
			time.Date(2023, 10, 1, 0, 0, 0, 0, time.Local),
			time.Date(2023, 10, 2, 12, 0, 0, 0, time.Local),
			[]time.Time{
				time.Date(2023, 10, 2, 12, 0, 0, 0, time.Local),
				time.Date(2023, 10, 3, 12, 0, 0, 0, time.Local),
				time.Date(2023, 10, 4, 12, 0, 0, 0, time.Local),
			},
			[]time.Time{
				time.Date(2023, 9, 29, 12, 0, 0, 0, time.Local),
				time.Date(2023, 9, 28, 12, 0, 0, 0, time.Local),
				time.Date(2023, 9, 27, 12, 0, 0, 0, time.Local),
			},
			false,
		},
		{
			"2023/9-10/* Mon,Fri,SUN 12:00:00",
			time.Date(2023, 10, 1, 0, 0, 0, 0, time.Local),
			time.Date(2023, 10, 1, 12, 0, 0, 0, time.Local),
			[]time.Time{
				time.Date(2023, 10, 1, 12, 0, 0, 0, time.Local),
				time.Date(2023, 10, 2, 12, 0, 0, 0, time.Local),
				time.Date(2023, 10, 6, 12, 0, 0, 0, time.Local),
			},
			[]time.Time{
				time.Date(2023, 9, 29, 12, 0, 0, 0, time.Local),
				time.Date(2023, 9, 25, 12, 0, 0, 0, time.Local),
				time.Date(2023, 9, 24, 12, 0, 0, 0, time.Local),
			},
			false,
		},
		{
			"2023/*/* 12,14,16,18,20,22:00:00",
			time.Date(2023, 10, 1, 0, 0, 0, 0, time.Local),
			time.Date(2023, 10, 1, 12, 0, 0, 0, time.Local),
			[]time.Time{
				time.Date(2023, 10, 1, 12, 0, 0, 0, time.Local),
				time.Date(2023, 10, 1, 14, 0, 0, 0, time.Local),
				time.Date(2023, 10, 1, 16, 0, 0, 0, time.Local),
			},
			[]time.Time{
				time.Date(2023, 9, 30, 22, 0, 0, 0, time.Local),
				time.Date(2023, 9, 30, 20, 0, 0, 0, time.Local),
				time.Date(2023, 9, 30, 18, 0, 0, 0, time.Local),
			},
			false,
		},
		{
			"2023/10/* * 12,23:00",
			time.Date(2023, 10, 8, 16, 0, 0, 0, time.Local),
			time.Date(2023, 10, 8, 23, 0, 0, 0, time.Local),
			[]time.Time{
				time.Date(2023, 10, 8, 23, 0, 0, 0, time.Local),
				time.Date(2023, 10, 9, 12, 0, 0, 0, time.Local),
				time.Date(2023, 10, 9, 23, 0, 0, 0, time.Local),
			},
			[]time.Time{
				time.Date(2023, 10, 8, 12, 0, 0, 0, time.Local),
				time.Date(2023, 10, 7, 23, 0, 0, 0, time.Local),
				time.Date(2023, 10, 7, 12, 0, 0, 0, time.Local),
			},
			false,
		},
		{
			"2023/10/* 17:00",
			time.Date(2023, 10, 5, 0, 0, 0, 0, time.Local),
			time.Date(2023, 10, 5, 17, 0, 0, 0, time.Local),
			[]time.Time{
				time.Date(2023, 10, 5, 17, 0, 0, 0, time.Local),
				time.Date(2023, 10, 6, 17, 0, 0, 0, time.Local),
				time.Date(2023, 10, 7, 17, 0, 0, 0, time.Local),
			},
			[]time.Time{
				time.Date(2023, 10, 4, 17, 0, 0, 0, time.Local),
				time.Date(2023, 10, 3, 17, 0, 0, 0, time.Local),
				time.Date(2023, 10, 2, 17, 0, 0, 0, time.Local),
			},
			false,
		},
		{
			"2023/*/* Mon,FRI 12:00",
			time.Date(2023, 10, 1, 0, 0, 0, 0, time.Local),
			time.Date(2023, 10, 2, 12, 0, 0, 0, time.Local),
			[]time.Time{
				time.Date(2023, 10, 2, 12, 0, 0, 0, time.Local),
				time.Date(2023, 10, 6, 12, 0, 0, 0, time.Local),
				time.Date(2023, 10, 9, 12, 0, 0, 0, time.Local),
			},
			[]time.Time{
				time.Date(2023, 9, 29, 12, 0, 0, 0, time.Local),
				time.Date(2023, 9, 25, 12, 0, 0, 0, time.Local),
				time.Date(2023, 9, 22, 12, 0, 0, 0, time.Local),
			},
			false,
		},
		{
			"2023/10/* MON 12:00:00",
			time.Date(2023, 10, 16, 12, 0, 0, 0, time.Local),
			time.Date(2023, 10, 16, 12, 0, 0, 0, time.Local),
			[]time.Time{
				time.Date(2023, 10, 23, 12, 0, 0, 0, time.Local),
				time.Date(2023, 10, 30, 12, 0, 0, 0, time.Local),
			},
			[]time.Time{
				time.Date(2023, 10, 9, 12, 0, 0, 0, time.Local),
				time.Date(2023, 10, 2, 12, 0, 0, 0, time.Local),
			},
			false,
		},
		{
			"2023/10/01 Mon 12:00:00:00",
			time.Date(2023, 10, 1, 0, 0, 0, 0, time.Local),
			time.Time{},
			[]time.Time{},
			[]time.Time{},
			true,
		},
		{
			"2023/10/01 Mon 12:00:xx",
			time.Date(2023, 10, 1, 0, 0, 0, 0, time.Local),
			time.Time{},
			[]time.Time{},
			[]time.Time{},
			true,
		},
		{
			"2023/10/01 Mon 12:xx:00",
			time.Date(2023, 10, 1, 0, 0, 0, 0, time.Local),
			time.Time{},
			[]time.Time{},
			[]time.Time{},
			true,
		},
		{
			"2023/10/01 Mon xx:00:00",
			time.Date(2023, 10, 1, 0, 0, 0, 0, time.Local),
			time.Time{},
			[]time.Time{},
			[]time.Time{},
			true,
		},
	}

	for _, test := range tests {
		t.Run(clearString(test.pattern), func(t *testing.T) {
			r := &RecurrentDatePattern{}
			err := r.ParseFromDatePattern(test.pattern)
			if (err != nil) != test.hasError {
				t.Errorf("ParseRecurrentDatePattern(%q) error = %v, wantErr %v", test.pattern, err, test.hasError)
				return
			}
			if err != nil {
				return
			}

			now := test.now
			for _, expectedNext := range test.expectedNext {
				next, err := r.Next(now)
				if err != nil {
					t.Fatalf("Next failed: %v", err)
				}
				if !next.Equal(expectedNext) {
					t.Errorf("Next(%v) = %v, want %v", now, next, expectedNext)
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

			first, err := r.First(test.now)
			if err != nil {
				t.Fatalf("Next failed: %v", err)
			}
			if !first.Equal(test.expectedFirst) {
				t.Errorf("Next(%v) = %v, want %v", now, first, test.expectedFirst)
			}
		})
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
			time.Date(2023, 10, 1, 0, 0, 0, 0, time.Local),
			time.Date(2023, 10, 1, 1, 30, 0, 0, time.Local),
			time.Date(2023, 9, 30, 22, 30, 0, 0, time.Local),
			false,
		},
		{
			"pattern(2023/10/* MON 12:00:00)",
			time.Date(2023, 10, 16, 0, 0, 0, 0, time.Local),
			time.Date(2023, 10, 16, 12, 0, 0, 0, time.Local),
			time.Date(2023, 10, 9, 12, 0, 0, 0, time.Local),
			false,
		},
		{
			"pattern(2023/*/* Mon-Fri 12:00:00)",
			time.Date(2023, 10, 1, 0, 0, 0, 0, time.Local),
			time.Date(2023, 10, 2, 12, 0, 0, 0, time.Local),
			time.Date(2023, 9, 29, 12, 0, 0, 0, time.Local),
			false,
		},
		{
			"pattern(2023/9-10/* Mon,Fri,SUN 12:00:00)",
			time.Date(2023, 10, 1, 0, 0, 0, 0, time.Local),
			time.Date(2023, 10, 1, 12, 0, 0, 0, time.Local),
			time.Date(2023, 9, 29, 12, 0, 0, 0, time.Local),
			false,
		},
		{
			"pattern(2023/*/* 12:00:00)",
			time.Date(2023, 10, 1, 0, 0, 0, 0, time.Local),
			time.Date(2023, 10, 1, 12, 0, 0, 0, time.Local),
			time.Date(2023, 9, 30, 12, 0, 0, 0, time.Local),
			false,
		},
		{
			"pattern(2023/10/01 * 12,23:00)",
			time.Date(2023, 10, 1, 16, 0, 0, 0, time.Local),
			time.Date(2023, 10, 1, 23, 0, 0, 0, time.Local),
			time.Date(2023, 10, 1, 12, 0, 0, 0, time.Local),
			false,
		},
		{
			"pattern(2023/10/* 17:00)",
			time.Date(2023, 10, 3, 0, 0, 0, 0, time.Local),
			time.Date(2023, 10, 3, 17, 0, 0, 0, time.Local),
			time.Date(2023, 10, 2, 17, 0, 0, 0, time.Local),
			false,
		},
		{
			"pattern(2023/*/* Mon,FRI 12:00)",
			time.Date(2023, 10, 1, 0, 0, 0, 0, time.Local),
			time.Date(2023, 10, 2, 12, 0, 0, 0, time.Local),
			time.Date(2023, 9, 29, 12, 0, 0, 0, time.Local),
			false,
		},
		{
			"pattern(12/25 12:00)",
			time.Date(2023, 10, 1, 0, 0, 0, 0, time.Local),
			time.Date(2023, 12, 25, 12, 0, 0, 0, time.Local),
			time.Date(2022, 12, 25, 12, 0, 0, 0, time.Local),
			false,
		},
		{
			"invalid(1h30m)",
			time.Date(2023, 10, 1, 0, 0, 0, 0, time.Local),
			time.Time{},
			time.Time{},
			true,
		},
		{
			"periodic(invalid)",
			time.Date(2023, 10, 1, 0, 0, 0, 0, time.Local),
			time.Time{},
			time.Time{},
			true,
		},
		{
			"pattern(invalid)",
			time.Date(2023, 10, 1, 0, 0, 0, 0, time.Local),
			time.Time{},
			time.Time{},
			true,
		},
		{
			"",
			time.Date(2023, 10, 1, 0, 0, 0, 0, time.Local),
			time.Time{},
			time.Time{},
			true,
		},
		{
			"()",
			time.Date(2023, 10, 1, 0, 0, 0, 0, time.Local),
			time.Time{},
			time.Time{},
			true,
		},
		{
			"periodic()",
			time.Date(2023, 10, 1, 0, 0, 0, 0, time.Local),
			time.Time{},
			time.Time{},
			true,
		},
		{
			"pattern()",
			time.Date(2023, 10, 1, 0, 0, 0, 0, time.Local),
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
		t.Run(clearString(test.pattern), func(t *testing.T) {
			recurrentDate, err := ParseRecurrentDate(test.pattern)
			if (err != nil) != test.hasError {
				t.Errorf("ParseRecurrentDate(%q) error = %v, wantErr %v", test.pattern, err, test.hasError)
				return
			}
			if err != nil {
				return
			}

			next, err := recurrentDate.Next(test.now)
			if err != nil {
				t.Fatalf("Next failed: %v", err)
			}
			if !next.Equal(test.expectedNext) {
				t.Errorf("Next(%v) = %v, want %v", test.now, next, test.expectedNext)
			}

			prev, err := recurrentDate.Prev(test.now)
			if err != nil {
				t.Fatalf("Prev failed: %v", err)
			}
			if !prev.Equal(test.expectedPrev) {
				t.Errorf("Prev(%v) = %v, want %v", test.now, prev, test.expectedPrev)
			}
		})
	}
}
