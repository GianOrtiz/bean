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

func Negative(m Money) Money {
	return FromCents(-1 * ToCents(m))
}

func Minus(m1 Money, m2 Money) Money {
	return FromCents(ToCents(m1) - ToCents(m2))
}

func Plus(m1 Money, m2 Money) Money {
	return FromCents(ToCents(m1) + ToCents(m2))
}
