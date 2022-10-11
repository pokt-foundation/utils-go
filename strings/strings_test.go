package strings

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient_ExactContains(t *testing.T) {
	c := require.New(t)

	array := []string{"dog", "cat", "shrew"}

	arrayContains := ExactContains(array, "shrew")
	c.True(arrayContains)

	arrayContains = ExactContains(array, "whale")
	c.False(arrayContains)
}
