package crypto

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCrypto_EncryptAES256(t *testing.T) {
	c := require.New(t)

	testCases := []struct {
		name          string
		inputString   string
		key           string
		expectedError bool
	}{
		{
			name:          "Should encryopt with a valid key",
			inputString:   "TEST_STRING",
			key:           "test_key_a1ed328e09ef0e0ca39b6b1",
			expectedError: false,
		},
		{
			name:          "Should fail to encrypt with an invalid key",
			inputString:   "TEST_STRING",
			key:           "invalid-key",
			expectedError: true,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			encryptedString, err := EncryptAES256(test.inputString, test.key)

			if test.expectedError {
				c.Error(err)
			} else {
				c.NoError(err)
				c.NotEmpty(encryptedString)
			}
		})
	}
}

func TestCrypto_EncryptAPIKey(t *testing.T) {
	c := require.New(t)

	testCases := []struct {
		name          string
		userData      UserData
		secretKey     string
		expectedError bool
	}{
		{
			name: "Should encrypt with a valid secret key",
			userData: UserData{
				UserID: "12345",
				Scope:  "free_scope",
			},
			secretKey:     "test_key_dd1e6bf2c074c09e222c14c",
			expectedError: false,
		},
		{
			name: "Should fail to encrypt with an invalid secret key",
			userData: UserData{
				UserID: "12345",
				Scope:  "free_scope",
			},
			secretKey:     "invalid-secret-key",
			expectedError: true,
		},
		{
			name: "Should fail to encrypt with an invalid key length",
			userData: UserData{
				UserID: "12345",
				Scope:  "free_scope",
			},
			secretKey:     "invalid_length_key",
			expectedError: true,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			keyUtil, err := NewKeyUtil(test.secretKey)
			if err != nil {
				if test.expectedError {
					c.Error(err)
					return
				}
				c.FailNow("Unexpected error:", err)
			}

			encryptedKey, err := keyUtil.EncryptAPIKey(test.userData)

			if test.expectedError {
				c.Error(err)
			} else {
				c.NoError(err)
				c.NotEmpty(encryptedKey)
			}
		})
	}
}

func TestCrypto_DecryptAPIKey(t *testing.T) {
	c := require.New(t)

	testCases := []struct {
		name          string
		encryptedKey  string
		secretKey     string
		expectedError bool
	}{
		{
			name:          "Should decrypt with a valid encrypted key and secret key",
			encryptedKey:  "p8HkIJD-hMSo5UltdqZvsCvEXZ0-L6CQk8BBnqGVX4clrvQQcRjommLT22b-TeIo0PHX", // pragma: allowlist secret
			secretKey:     "test_key_dd1e6bf2c074c09e222c14c",
			expectedError: false,
		},
		{
			name:          "Should fail to decrypt with an invalid encrypted key",
			encryptedKey:  "invalid-encrypted-key",
			secretKey:     "test_key_dd1e6bf2c074c09e222c14c",
			expectedError: true,
		},
		{
			name:          "Should fail to decrypt with an invalid secret key",
			encryptedKey:  "test_key_4DWjg4e1n6Y96bmHRHrhK4rl0CmMfRgDzQ2HyY",
			secretKey:     "invalid-secret-key",
			expectedError: true,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			keyUtil, err := NewKeyUtil(test.secretKey)
			if err != nil {
				if test.expectedError {
					c.Error(err)
					return
				}
				c.FailNow("Unexpected error:", err)
			}

			userData, err := keyUtil.DecryptAPIKey(test.encryptedKey)

			if test.expectedError {
				c.Error(err)
			} else {
				c.NoError(err)
				c.NotEmpty(userData.UserID)
				c.NotEmpty(userData.Scope)
			}
		})
	}
}
