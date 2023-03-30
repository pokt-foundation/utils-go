// Package id handles all ID related funcs
package id

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateID generates a random ID with the specified length
func GenerateID(length int) string {
	bytesLength := length / 2
	if length%2 != 0 {
		bytesLength++
	}
	bytes := make([]byte, bytesLength)
	_, _ = rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
