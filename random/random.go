package random

import (
	"crypto/rand"
	"encoding/hex"
)

// HexString returns a random hex string of n length
func HexString(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
