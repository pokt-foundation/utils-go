package numbers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient_RoundFloat(t *testing.T) {
	c := require.New(t)

	float := 0.12774863846293874093274

	roundedFloat := RoundFloat(float, 5)
	c.Equal(0.12775, roundedFloat)

	roundedFloat = RoundFloat(float, 7)
	c.Equal(0.1277486, roundedFloat)

	roundedFloat = RoundFloat(float, 2)
	c.Equal(0.13, roundedFloat)
}

func TestClient_RoundDownFloat(t *testing.T) {
	c := require.New(t)

	float := 0.12774863846293874093274

	roundedDownFloat := RoundDownFloat(float, 5)
	c.Equal(0.12774, roundedDownFloat)

	roundedDownFloat = RoundDownFloat(float, 7)
	c.Equal(0.1277486, roundedDownFloat)

	roundedDownFloat = RoundDownFloat(float, 2)
	c.Equal(0.12, roundedDownFloat)
}

func TestClient_RoundUpFloat(t *testing.T) {
	c := require.New(t)

	float := 0.12774863846293874093274

	roundedUpFloat := RoundUpFloat(float, 5)
	c.Equal(0.12775, roundedUpFloat)

	roundedUpFloat = RoundUpFloat(float, 7)
	c.Equal(0.1277487, roundedUpFloat)

	roundedUpFloat = RoundUpFloat(float, 2)
	c.Equal(0.13, roundedUpFloat)
}
