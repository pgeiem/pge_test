package engine

import (
	"fmt"
	"math"
)

// Amount represents the amount of money represented a fixed point value in micro units

type Amount float64 //TODO replace by int64, and update min/max values

const AmountMax = Amount(math.MaxFloat64)

// String returns the string representation of the amount
func (a Amount) String() string {

	return fmt.Sprintf("%.2f", a)
}
