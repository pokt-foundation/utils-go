// Package numbers includes all funcs for manipulating and examining numbers
package numbers

import "math"

// RoundFloat rounds a float to a specified precision
func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
