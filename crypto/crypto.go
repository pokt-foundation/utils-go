// Package crypto handles all encryption related funcs
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

// EncryptAES256 encrypts a string using AES256 encryption and returns the encrypted string
func EncryptAES256(stringToEncrypt string, keyString string) (string, error) {
	plaintext := []byte(stringToEncrypt)
	plainKey := []byte(keyString)

	block, err := aes.NewCipher(plainKey)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	encryptedString := fmt.Sprintf("%x", ciphertext)

	return encryptedString, nil
}
