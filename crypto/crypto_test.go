package crypto

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient_EncryptAES256(t *testing.T) {
	c := require.New(t)

	testKey := "test_key_a1ed328e09ef0e0ca39b6b1"

	encryptedString, err := EncryptAES256("TEST_STRING", testKey)
	c.NoError(err)
	c.NotEmpty(encryptedString)
	c.Len(encryptedString, 78)

	_, err = EncryptAES256("TEST_STRING", "invalid-key")
	c.Error(err)
}
