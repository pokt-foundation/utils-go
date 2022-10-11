package crypto

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient_Encrypt(t *testing.T) {
	c := require.New(t)

	testKey := "test_key_a1ed328e09ef0e0ca39b6b1"

	encryptedString, err := Encrypt("TEST_STRING", testKey)
	c.NoError(err)
	c.NotEmpty(encryptedString)
	c.Len(encryptedString, 78)

	_, err = Encrypt("TEST_STRING", "invalid-key")
	c.Error(err)
}
