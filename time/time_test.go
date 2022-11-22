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
