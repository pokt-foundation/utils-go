// Package crypto handles all encryption related funcs
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"
)

type (
	KeyUtil struct {
		secretKey string
	}
	UserData struct {
		UserID, Scope string
	}
	IKeyUtil interface {
		EncryptAPIKey(userData UserData) (string, error)
		DecryptAPIKey(encryptedKey string) (UserData, error)
	}
)

func NewKeyUtil(secretKey string) (IKeyUtil, error) {
	keyLength := len(secretKey)

	if keyLength != 16 && keyLength != 24 && keyLength != 32 {
		return nil, errors.New("invalid secret key length: must be 16, 24, or 32 bytes")
	}

	return &KeyUtil{secretKey: secretKey}, nil
}

// EncryptAPIKey encrypts the user data and returns the encrypted API key.
func (k *KeyUtil) EncryptAPIKey(userData UserData) (string, error) {
	combined := fmt.Sprintf("%s:%s", userData.UserID, userData.Scope)

	block, err := aes.NewCipher([]byte(k.secretKey))
	if err != nil {
		return "", err
	}

	nonce := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	stream := cipher.NewCTR(block, nonce)

	ciphertext := make([]byte, len(combined))
	stream.XORKeyStream(ciphertext, []byte(combined))

	result := make([]byte, len(nonce)+len(ciphertext))
	copy(result[:len(nonce)], nonce)
	copy(result[len(nonce):], ciphertext)

	return base64.URLEncoding.EncodeToString(result), nil
}

// DecryptAPIKey decrypts the encrypted API key and returns the user data.
func (k *KeyUtil) DecryptAPIKey(encryptedKey string) (UserData, error) {
	data, err := base64.URLEncoding.DecodeString(encryptedKey)
	if err != nil {
		return UserData{}, err
	}

	if len(data) < 16 {
		return UserData{}, errors.New("encrypted data too short")
	}

	nonce, ciphertext := data[:16], data[16:]

	block, err := aes.NewCipher([]byte(k.secretKey))
	if err != nil {
		return UserData{}, err
	}

	stream := cipher.NewCTR(block, nonce)

	plaintext := make([]byte, len(ciphertext))
	stream.XORKeyStream(plaintext, ciphertext)

	parts := strings.Split(string(plaintext), ":")
	if len(parts) != 2 {
		return UserData{}, errors.New("decrypted data has an invalid format")
	}

	userData := UserData{
		UserID: parts[0],
		Scope:  parts[1],
	}

	return userData, nil
}

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
