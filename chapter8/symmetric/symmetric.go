package symmetric

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
)

// EncryptData encrypts data using AES with the given key
func EncryptData(data, key []byte) ([]byte, error) {
	if err := validateKey(key); err != nil {
		return make([]byte, 0), err
	}

	c, err := aes.NewCipher(key)
	if err != nil {
		return make([]byte, 0), err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return make([]byte, 0), err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return make([]byte, 0), err
	}

	return gcm.Seal(nonce, nonce, data, nil), nil
}

// DecryptData decrypts data with given key
func DecryptData(data, key []byte) ([]byte, error) {
	if err := validateKey(key); err != nil {
		return make([]byte, 0), err
	}

	c, err := aes.NewCipher(key)
	if err != nil {
		return make([]byte, 0), err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return make([]byte, 0), err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return make([]byte, 0), errors.New("ciphertext too short")
	}

	nonce, cipherText := data[:nonceSize], data[nonceSize:]
	return gcm.Open(nil, nonce, cipherText, nil)
}

func validateKey(key []byte) error {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return errors.New("Invalid key length, keys should be 16, 24, or 32 bytes in length")
	}

	return nil
}
