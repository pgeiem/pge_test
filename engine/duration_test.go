package engine

import (
	"testing"
	"time"
)

func TestParseDuration(t *testing.T) {
	tests := []struct {
		input    string
		expected time.Duration
		hasError bool
	}{
		{"1w", time.Duration(7 * 24 * time.Hour), false},
		{"2d", time.Duration(2 * 24 * time.Hour), false},
		{"3h", time.Duration(3 * time.Hour), false},
		{"4m", time.Duration(4 * time.Minute), false},
		{"5s", time.Duration(5 * time.Second), false},
		{"1w2d3h4m5s", time.Duration(7*24*time.Hour + 2*24*time.Hour + 3*time.Hour + 4*time.Minute + 5*time.Second), false},
		{"2d3h4m5s", time.Duration(2*24*time.Hour + 3*time.Hour + 4*time.Minute + 5*time.Second), false},
		{"3h4m5s", time.Duration(3*time.Hour + 4*time.Minute + 5*time.Second), false},
		{"4m5s", time.Duration(4*time.Minute + 5*time.Second), false},
		{"5s", time.Duration(5 * time.Second), false},
		{"1w2h", time.Duration(7*24*time.Hour + 2*time.Hour), false}, // Valid order
		{"2d1w", 0, true}, // Invalid order
		{"3h2d", 0, true}, // Invalid order
		{"4m3h", 0, true}, // Invalid order
		{"5s4m", 0, true}, // Invalid order
		{"", 0, true},     // Empty string
		{"1x", 0, true},   // Unknown unit
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result, err := ParseDuration(test.input)
			if (err != nil) != test.hasError {
				t.Errorf("ParseDuration(%q) error = %v, wantErr %v", test.input, err, test.hasError)
				return
			}
			if result != test.expected {
				t.Errorf("ParseDuration(%q) = %v, want %v", test.input, result, test.expected)
			}
		})
	}
}
