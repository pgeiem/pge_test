package engine

import (
	"testing"
)

func TestParseAmount(t *testing.T) {
	tests := []struct {
		input    string
		expected Amount
		hasError bool
	}{
		{"0", Amount(0), false},
		{"2.5", Amount(2500000), false},
		{"1.50", Amount(1500000), false},
		{"5", Amount(5000000), false},
		{"123.5678", Amount(123567800), false},
		{"2147.48", Amount(2147480000), false},
		{"2147.50", Amount(0), true},
		{"-2.5", Amount(-2500000), false},
		{"-1.50", Amount(-1500000), false},
		{"-5", Amount(-5000000), false},
		{"-2147.48", Amount(-2147480000), false},
		{"-2147.50", Amount(0), true},
		{"+2.5", Amount(2500000), false},
		{"+1.50", Amount(1500000), false},
		{"+5", Amount(5000000), false},
		{"+123.5678", Amount(123567800), false},
		{"+2147.48", Amount(2147480000), false},
		{"+2147.50", Amount(0), true},
		{"1.123456", Amount(1123456), false},
		{"+123.123456789012345", Amount(0), true}, // More than 6 decimal places
		{"1.1234567", Amount(0), true},            // More than 6 decimal places
		{"abc", Amount(0), true},                  // Invalid format
		{"", Amount(0), true},                     // Empty string
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result, err := ParseAmount(test.input)
			if (err != nil) != test.hasError {
				t.Errorf("ParseAmount(%s) error = %v, expected error = %v", test.input, err, test.hasError)
			}
			if result != test.expected {
				t.Errorf("ParseAmount(%s) = %v, expected %v", test.input, result, test.expected)
			}
		})
	}
}
