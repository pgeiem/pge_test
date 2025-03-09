package engine

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
)

// Amount represents the amount of money represented a fixed point value in micro units
// The amount is stored as an unsigned integer to avoid floating point rounding errors
// For example 1.50€ is stored as 1500000.
// Maximal amount is 2147.48 and minimal amount is -2147.48
type Amount int32 //TODO replace by int64, and update min/max values

const (
	// AmountZero represents the zero amount
	AmountZero Amount = 0
	// AmountUnit represents the one amount
	AmountUnit Amount = 1000000
	// AmountMax represents the maximal amount
	AmountMax Amount = 2147480000
)

var amountRegex = regexp.MustCompile(`^([+-]?\d+)(?:.(\d{0,6}))?$`)

// ParseAmount parses a string representing an amount of money
func ParseAmount(s string) (Amount, error) {
	var a int64
	var err error
	// Parse the amount
	matches := amountRegex.FindStringSubmatch(s)
	if matches == nil || s == "" {
		return AmountZero, fmt.Errorf("error while parsing %s amount, invalid format", s)
	}
	// Convert the integer part
	if matches[1] != "" {
		a, err = strconv.ParseInt(matches[1], 10, 32)
		a *= int64(AmountUnit)
		if err != nil {
			return AmountZero, fmt.Errorf("error while parsing %s amount, invalid format (integer part)", s)
		}
	}
	// Convert the fractional part
	if matches[2] != "" {
		multipliers := []Amount{AmountUnit, AmountUnit / 10, AmountUnit / 100, AmountUnit / 1000, AmountUnit / 10000, AmountUnit / 100000, AmountUnit / 1000000}
		afrac, err := strconv.ParseInt(matches[2], 10, 32)
		if err != nil {
			return AmountZero, fmt.Errorf("error while parsing %s amount, invalid format (fractional part)", s)
		}
		afrac *= int64(multipliers[len(matches[2])])
		if a >= 0 {
			a += afrac
		} else {
			a -= afrac
		}
	}
	//Check for overflow
	if a >= math.MaxInt32 || a <= math.MinInt32 {
		return AmountZero, fmt.Errorf("error while parsing %s amount, too large", s)
	}
	return Amount(a), nil
}

func MustParseAmount(s string) Amount {
	amount, err := ParseAmount(s)
	if err != nil {
		panic(err)
	}
	return amount
}

// MarshalText implements the encoding.TextMarshaler interface
func (a *Amount) UnmarshalText(text []byte) error {
	amount, err := ParseAmount(string(text))
	if err != nil {
		return err
	}
	*a = amount
	return nil
}

func (a *Amount) MarshalText() ([]byte, error) {
	return []byte(a.String()), nil
}

// String returns the string representation of the amount
func (a Amount) String() string {
	if a < 0 {
		return fmt.Sprintf("-%d.%03d", -a/AmountUnit, -a%AmountUnit)
	}
	return fmt.Sprintf("%d.%03d", a/AmountUnit, a%AmountUnit)
}
