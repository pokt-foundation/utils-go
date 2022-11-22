// Package time includes all funcs for dealing with times and dates
package time

import "time"

// GetFirstDayOfMonth gets the first day of the month for any given time.Time value
func GetFirstDayOfMonth(day time.Time) time.Time {
	currentYear, currentMonth, _ := day.Date()
	currentLocation := day.Location()
	return time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
}
