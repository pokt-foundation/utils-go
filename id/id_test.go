package id

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateID(t *testing.T) {
	c := require.New(t)

	testCases := []struct {
		name   string
		length int
	}{
		{
			name:   "Should generate an ID of length 6",
			length: 6,
		},
		{
			name:   "Should generate another ID of length 6",
			length: 6,
		},
		{
			name:   "Should generate an ID of length 7",
			length: 7,
		},
		{
			name:   "Should generate another ID of length 7",
			length: 7,
		},
	}

	ids := make(map[string]bool)

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			id := GenerateID(test.length)

			c.Len(id, (test.length+1)/2*2)
			c.NotContains(ids, id)

			ids[id] = true
		})
	}
}
