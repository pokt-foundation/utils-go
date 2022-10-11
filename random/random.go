// Package random includes all funcs for generating randomness
package random

import (
	"crypto/rand"
	"encoding/hex"
)

// HexString returns a random hex string of n length
func HexString(n int) (string, error) {
	stringLength := n / 2

	bytes := make([]byte, stringLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	randomString := hex.EncodeToString(bytes)

	return randomString, nil
}
