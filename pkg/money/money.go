package money

import (
	"math"
)

type Money int

func ToCents(m Money) int {
	return int(m)
}

func FromFloat(number float64) Money {
	roundedNumber := math.Round(number*100) / 100
	return Money(int(roundedNumber))
}

func FromCents(cents int) Money {
	return Money(cents)
}
