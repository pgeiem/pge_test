package engine

import (
	"testing"
	"time"
)

func mustParseRecurrentDate(pattern string) RecurrentDate {
	r, err := ParseRecurrentDate(pattern)
	if err != nil {
		panic(err)
	}
	return r
}

func TestQuota_Update(t *testing.T) {
	tests := []struct {
		name             string
		now              time.Time
		matchingRules    []MatchingRule
		periodicityRule  RecurrentDate
		history          []AssignedRight
		expectedDuration time.Duration
		expectedCounter  int
		expectedError    bool
	}{
		{
			name: "Single area with free and paying duration",
			now:  time.Date(2023, 10, 10, 12, 30, 0, 0, time.UTC),
			matchingRules: []MatchingRule{
				{ParkingAreaPattern: "area1", DurationTypePattern: "free"},
			},
			periodicityRule: mustParseRecurrentDate("duration(1d)"),
			history: []AssignedRight{
				{
					ParkingArea: []string{"area1"},
					Details: []DurationDetail{
						{Type: FreeDuration, Start: time.Date(2023, 10, 10, 4, 0, 0, 0, time.UTC), Duration: 1 * time.Hour},
						{Type: PayingDuration, Start: time.Date(2023, 10, 9, 12, 30, 0, 0, time.UTC), Duration: 2 * time.Hour},
						{Type: FreeDuration, Start: time.Date(2023, 10, 9, 12, 0, 0, 0, time.UTC), Duration: 4 * time.Hour},
					},
				},
			},
			expectedDuration: 1 * time.Hour,
			expectedCounter:  1,
			expectedError:    false,
		},
		{
			name: "Empty History",
			now:  time.Date(2023, 10, 10, 12, 30, 0, 0, time.UTC),
			matchingRules: []MatchingRule{
				{ParkingAreaPattern: "area1", DurationTypePattern: "free"},
			},
			periodicityRule:  mustParseRecurrentDate("duration(1d)"),
			expectedDuration: 0,
			expectedCounter:  0,
			expectedError:    false,
		},
		{
			name: "Multiple areas with free duration",
			now:  time.Date(2023, 10, 4, 0, 0, 0, 0, time.UTC),
			matchingRules: []MatchingRule{
				{ParkingAreaPattern: "area*", DurationTypePattern: "free"},
			},
			periodicityRule: mustParseRecurrentDate("duration(2d)"),
			history: []AssignedRight{
				{
					ParkingArea: []string{"area1"},
					Start:       time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC),
					Details: []DurationDetail{
						{Type: FreeDuration, Duration: 2 * time.Hour},
					},
				},
				{
					ParkingArea: []string{"area2"},
					Start:       time.Date(2023, 10, 3, 0, 0, 0, 0, time.UTC),
					Details: []DurationDetail{
						{Type: FreeDuration, Duration: 3 * time.Hour},
					},
				},
				{
					ParkingArea: []string{"area3"},
					Start:       time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
					Details: []DurationDetail{
						{Type: FreeDuration, Duration: 3 * time.Hour},
					},
				},
			},
			expectedDuration: 5 * time.Hour,
			expectedCounter:  3,
			expectedError:    false,
		},
		{
			name:            "No matching rules",
			now:             time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC),
			matchingRules:   []MatchingRule{},
			periodicityRule: mustParseRecurrentDate("duration(1d)"),
			history: []AssignedRight{
				{
					ParkingArea: []string{"area1"},
					Start:       time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
					Details: []DurationDetail{
						{Type: PayingDuration, Duration: 1 * time.Hour},
						{Type: FreeDuration, Duration: 2 * time.Hour},
					},
				},
			},
			expectedDuration: 2 * time.Hour,
			expectedCounter:  1,
			expectedError:    false,
		},
		{
			name: "No start defined",
			now:  time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			matchingRules: []MatchingRule{
				{ParkingAreaPattern: "area*", DurationTypePattern: "free"},
			},
			periodicityRule: mustParseRecurrentDate("duration(1d)"),
			history: []AssignedRight{
				{
					ParkingArea: []string{"area1"},
					Details: []DurationDetail{
						{Type: PayingDuration, Duration: 1 * time.Hour},
						{Type: FreeDuration, Duration: 2 * time.Hour},
					},
				},
			},
			expectedDuration: 0,
			expectedCounter:  1,
			expectedError:    false,
		},
		{
			name: "Multiple areas with mixed types",
			now:  time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC),
			matchingRules: []MatchingRule{
				{ParkingAreaPattern: "area*", DurationTypePattern: "free"},
			},
			periodicityRule: mustParseRecurrentDate("duration(1d)"),
			history: []AssignedRight{
				{
					ParkingArea: []string{"area1"},
					Details: []DurationDetail{
						{Type: FreeDuration, Start: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC), Duration: 1 * time.Hour},
						{Type: PayingDuration, Start: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC), Duration: 2 * time.Hour},
					},
				},
				{
					ParkingArea: []string{"area2"},
					Details: []DurationDetail{
						{Type: FreeDuration, Start: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC), Duration: 4 * time.Hour},
						{Type: PayingDuration, Start: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC), Duration: 8 * time.Hour},
						{Type: FreeDuration, Start: time.Date(2023, 9, 30, 0, 0, 0, 0, time.UTC), Duration: 16 * time.Hour},
					},
				},
			},
			expectedDuration: 5 * time.Hour,
			expectedCounter:  2,
			expectedError:    false,
		},
		{
			name: "Single area with multiple free durations",
			now:  time.Date(2023, 10, 15, 0, 0, 0, 0, time.UTC),
			matchingRules: []MatchingRule{
				{ParkingAreaPattern: "area1", DurationTypePattern: "free"},
			},
			periodicityRule: mustParseRecurrentDate("pattern(*/*/* MON 12:00:00)"),
			history: []AssignedRight{
				{
					ParkingArea: []string{"area1"},
					Details: []DurationDetail{
						{Type: FreeDuration, Start: time.Date(2023, 10, 13, 0, 0, 0, 0, time.UTC), Duration: 1 * time.Hour},
						{Type: FreeDuration, Start: time.Date(2023, 10, 9, 12, 0, 0, 0, time.UTC), Duration: 2 * time.Hour},
						{Type: FreeDuration, Start: time.Date(2023, 10, 9, 11, 0, 0, 0, time.UTC), Duration: 4 * time.Hour},
						{Type: FreeDuration, Start: time.Date(2023, 10, 8, 0, 0, 0, 0, time.UTC), Duration: 8 * time.Hour},
					},
				},
			},
			expectedDuration: 3 * time.Hour,
			expectedCounter:  1,
			expectedError:    false,
		},
		{
			name: "Different area pattern",
			now:  time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			matchingRules: []MatchingRule{
				{ParkingAreaPattern: "area2", DurationTypePattern: "free"},
			},
			periodicityRule: mustParseRecurrentDate("duration(1d)"),
			history: []AssignedRight{
				{
					ParkingArea: []string{"area1"},
					Details: []DurationDetail{
						{Type: FreeDuration, Start: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC), Duration: 1 * time.Hour},
					},
				},
				{
					ParkingArea: []string{"area3", "area2"},
					Details: []DurationDetail{
						{Type: FreeDuration, Start: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC), Duration: 2 * time.Hour},
					},
				},
			},
			expectedDuration: 2 * time.Hour,
			expectedCounter:  1,
			expectedError:    false,
		},
		{
			name: "Glob pattern for Type",
			now:  time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			matchingRules: []MatchingRule{
				{ParkingAreaPattern: "area1", DurationTypePattern: "free*"},
			},
			periodicityRule: mustParseRecurrentDate("duration(1d)"),
			history: []AssignedRight{
				{
					ParkingArea: []string{"area1"},
					Details: []DurationDetail{
						{Type: FreeDuration, Start: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC), Duration: 1 * time.Hour},
						{Type: FreeDuration, Start: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC), Duration: 2 * time.Hour},
						{Type: NonPayingDuration, Start: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC), Duration: 4 * time.Hour},
					},
				},
			},
			expectedDuration: 3 * time.Hour,
			expectedCounter:  1,
			expectedError:    false,
		},
		{
			name: "Multiple matching rules",
			now:  time.Date(2023, 10, 1, 20, 0, 0, 0, time.UTC),
			matchingRules: []MatchingRule{
				{ParkingAreaPattern: "area1", DurationTypePattern: "free"},
				{ParkingAreaPattern: "area2", DurationTypePattern: "paying"},
			},
			periodicityRule: mustParseRecurrentDate("pattern(*/*/* 12:00:00)"),
			history: []AssignedRight{
				{
					ParkingArea: []string{"area1"},
					Details: []DurationDetail{
						{Type: FreeDuration, Start: time.Date(2023, 10, 1, 13, 0, 0, 0, time.UTC), Duration: 1 * time.Hour},
						{Type: PayingDuration, Start: time.Date(2023, 10, 1, 14, 0, 0, 0, time.UTC), Duration: 2 * time.Hour},
					},
				},
				{
					ParkingArea: []string{"area4", "area2"},
					Start:       time.Date(2023, 10, 1, 15, 0, 0, 0, time.UTC),
					Details: []DurationDetail{
						{Type: PayingDuration, Duration: 4 * time.Hour},
						{Type: FreeDuration, Duration: 8 * time.Hour},
						{Type: PayingDuration, Start: time.Date(2023, 10, 1, 10, 0, 0, 0, time.UTC), Duration: 16 * time.Hour},
					},
				},
			},
			expectedDuration: 5 * time.Hour,
			expectedCounter:  2,
			expectedError:    false,
		},
		{
			name: "Invalid area pattern",
			now:  time.Date(2023, 10, 1, 20, 0, 0, 0, time.UTC),
			matchingRules: []MatchingRule{
				{ParkingAreaPattern: "]", DurationTypePattern: "free"},
			},
			periodicityRule: mustParseRecurrentDate("pattern(*/*/* 12:00:00)"),
			history: []AssignedRight{
				{
					ParkingArea: []string{"area1"},
					Details:     []DurationDetail{},
				},
			},
			expectedDuration: 0,
			expectedCounter:  0,
			expectedError:    false,
		},
		{
			name: "Invalid type pattern",
			now:  time.Date(2023, 10, 1, 20, 0, 0, 0, time.UTC),
			matchingRules: []MatchingRule{
				{ParkingAreaPattern: "*", DurationTypePattern: "]"},
			},
			periodicityRule: mustParseRecurrentDate("pattern(*/*/* 12:00:00)"),
			history: []AssignedRight{
				{
					ParkingArea: []string{"area1"},
					Details:     []DurationDetail{},
				},
			},
			expectedDuration: 0,
			expectedCounter:  1,
			expectedError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Test the duration quota
			durationQuota := NewDurationQuota(time.Duration(0), tt.periodicityRule, tt.matchingRules)
			err := durationQuota.Update(tt.now, tt.history)
			if (err != nil) != tt.expectedError {
				t.Fatalf("expected error %v, got %v", tt.expectedError, err)
			}

			if durationQuota.Used() != tt.expectedDuration {
				t.Errorf("expected duration %v, got %v", tt.expectedDuration, durationQuota.Used())
			}

			// Tests the counter quota
			counterQuota := NewCounterQuota(0, tt.periodicityRule, tt.matchingRules)
			err = counterQuota.Update(tt.now, tt.history)
			if tt.expectedError {
				t.Fatalf("expected error %v, got %v", tt.expectedError, err)
			}

			if (err == nil) && counterQuota.Used() != tt.expectedCounter {
				t.Errorf("expected counter %v, got %v", tt.expectedCounter, counterQuota.Used())
			}
		})
	}
}
