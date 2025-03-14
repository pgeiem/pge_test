package engine

import (
	"fmt"
	"math"
)

// Amount represents the amount of money represented a floating point value in unit
// The currency has not importane here, it is just a number

type Amount float64

const AmountMax = Amount(math.MaxFloat64)

// String returns the string representation of the amount
func (a Amount) String() string {
	return fmt.Sprintf("%.2f", a)
}

// Reduce float precision to 6 decimal places, usefull when exporting to json
func (a Amount) Simplify() Amount {
	return Amount(math.Round(float64(a)*1000000) / 1000000)
}
