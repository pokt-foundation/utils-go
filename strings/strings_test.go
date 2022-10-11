package strings

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient_Contains(t *testing.T) {
	c := require.New(t)

	array := []string{"dog", "cat", "shrew"}

	arrayContains := Contains(array, "shrew")
	c.True(arrayContains)

	arrayContains = Contains(array, "cat")
	c.True(arrayContains)

	arrayContains = Contains(array, "whale")
	c.False(arrayContains)
}
