// Package numbers includes all funcs for manipulating and examining numbers
package numbers

import "math"

// RoundFloat rounds a float to a specified precision
func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

// RoundDownFloat rounds a float down to a specified precision
func RoundDownFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Floor(val*ratio) / ratio
}

// RoundUpFloat rounds a float up to a specified precision
func RoundUpFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Ceil(val*ratio) / ratio
}
