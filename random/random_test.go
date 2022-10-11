package random

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient_HexString(t *testing.T) {
	c := require.New(t)

	length := 32

	randomHex, err := HexString(length)
	c.NoError(err)
	c.Len(randomHex, length)
}
