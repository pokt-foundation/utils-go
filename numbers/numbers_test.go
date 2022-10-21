package numbers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient_RoundFloat(t *testing.T) {
	c := require.New(t)

	float := 0.12374863846293874093274

	roundedFloat := RoundFloat(float, 5)
	c.Equal(0.12375, roundedFloat)

	roundedFloat = RoundFloat(float, 7)
	c.Equal(0.1237486, roundedFloat)
}
