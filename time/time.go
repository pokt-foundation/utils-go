// Package time includes all funcs for dealing with times and dates
package time

import "time"

// GetFirstDayOfMonth gets the first day of the month for any given time.Time value
func GetFirstDayOfMonth(day time.Time) time.Time {
	currentYear, currentMonth, _ := day.Date()
	currentLocation := day.Location()
	return time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
}

// DiffMonths calculates the number of months between two dates
func DiffMonths(start, end time.Time) int {
	diffYears := end.Year() - start.Year()
	if diffYears == 0 {
		return int(end.Month() - start.Month())
	}

	if diffYears == 1 {
		return MonthsTillEndOfYear(start) + int(end.Month())
	}

	yearsInMonths := (end.Year() - start.Year() - 1) * 12
	return yearsInMonths + MonthsTillEndOfYear(start) + int(end.Month())
}

// MonthsTillEndOfYear calculates the number of months until the end of the year
func MonthsTillEndOfYear(month time.Time) int {
	return int(12 - month.Month())
}
