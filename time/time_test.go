package time

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var dayLayout = "2006-01-02"

func TestTime_GetFirstDayOfMonth(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		date, expected string
	}{
		{date: "2022-11-12", expected: "2022-11-01"},
		{date: "2022-11-01", expected: "2022-11-01"},
		{date: "2012-12-21", expected: "2012-12-01"},
		{date: "2032-01-31", expected: "2032-01-01"},
	}

	for _, test := range tests {
		t.Run(test.date, func(t *testing.T) {
			date, _ := time.Parse(dayLayout, test.date)

			firstDayOfMonth := GetFirstDayOfMonth(date)
			c.Equal(test.expected, firstDayOfMonth.Format(dayLayout))
		})
	}
}

func TestTime_DiffMonths(t *testing.T) {
	c := require.New(t)

	testDate, _ := time.Parse(dayLayout, "2022-11-12")

	tests := []struct {
		name           string
		start, end     time.Time
		expectedMonths int
	}{
		{
			name:           "one month",
			start:          testDate,
			end:            testDate.AddDate(0, 1, 0),
			expectedMonths: 1,
		},
		{
			name:           "one year",
			start:          testDate,
			end:            testDate.AddDate(1, 0, 0),
			expectedMonths: 12,
		},
		{
			name:           "sixty days",
			start:          testDate,
			end:            testDate.AddDate(0, 0, 60),
			expectedMonths: 2,
		},
		{
			name:           "one year, two months, fourteen days",
			start:          testDate,
			end:            testDate.AddDate(1, 2, 14),
			expectedMonths: 14,
		},
		{
			name:           "two months earlier",
			start:          testDate,
			end:            testDate.AddDate(0, -2, 0),
			expectedMonths: -2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			months := DiffMonths(test.start, test.end)

			c.Equal(test.expectedMonths, months)
		})
	}
}

func TestTime_MonthsTillEndOfYear(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name, date     string
		expectedMonths int
	}{
		{name: "January", date: "2022-01-01", expectedMonths: 11},
		{name: "June", date: "2022-06-26", expectedMonths: 6},
		{name: "December", date: "2022-12-25", expectedMonths: 0},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			date, _ := time.Parse(dayLayout, test.date)
			months := MonthsTillEndOfYear(date)

			c.Equal(test.expectedMonths, months)
		})
	}
}
