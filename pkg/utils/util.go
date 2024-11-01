package utils

import "math"

func RoundToInt(val float64, precision int) int32 {
	p := math.Pow10(precision)
	return int32(math.Floor(val*p+0.5) / p)
}

func RoundToDecimal(val float64, precision int) float64 {
	p := math.Pow10(precision)
	return math.Floor(val*p+0.5) / p
}
