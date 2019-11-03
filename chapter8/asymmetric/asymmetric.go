package asymmetric

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"

	"github.com/dikaeinstein/go-microservice/chapter8/symmetric"

	"github.com/dikaeinstein/go-microservice/chapter8/crypto"
)

var rsaPublic *rsa.PublicKey
var rsaPrivate *rsa.PrivateKey

func init() {
	var err error
	rsaPrivate, err = crypto.UnmarshalRSAPrivateKeyFromFile("../keys/sample_key.priv")
	if err != nil {
		log.Fatal("Unable to read private key", err)
	}

	rsaPublic, err = crypto.UnmarshalRSAPublicKeyFromFile("../keys/sample_key.pub")
	if err != nil {
		log.Fatal("Unable to read public key", err)
	}
}

// EncryptDataWithPublicKey encrypts the given data with the public key
func EncryptDataWithPublicKey(data []byte) ([]byte, error) {
	return rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaPublic, data, nil)
}

// EncryptMessageWithPublicKey encrypts the given string and returns the encrypted
// result base64 encoded
func EncryptMessageWithPublicKey(message string) (string, error) {
	modulus := rsaPublic.N.BitLen() / 8
	hashLength := 256 / 4
	maxLength := modulus - (hashLength * 2) - 2

	if len(message) > maxLength {
		return "", fmt.Errorf("The maximum message size must not exceed: %d", maxLength)
	}

	data, err := EncryptDataWithPublicKey([]byte(message))
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(data), nil
}

// DecryptDataWithPrivateKey decrypts the given data with the private key
func DecryptDataWithPrivateKey(data []byte) ([]byte, error) {
	return rsa.DecryptOAEP(sha256.New(), rand.Reader, rsaPrivate, data, nil)
}

// DecryptMessageWithPrivateKey decrypts the given base64 encoded ciphertext with
// the private key and returns plain text
func DecryptMessageWithPrivateKey(message string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return "", err
	}

	data, err = DecryptDataWithPrivateKey(data)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// EncryptLargeMessageWithPublicKey encrypts the given message by randomly generating
// a cipher. Returns the ciphertext for the given message base64 encoded and the key
func EncryptLargeMessageWithPublicKey(message string) (cipherText, cipherKey string, err error) {
	key, err := crypto.GenerateRandomString(16)
	if err != nil {
		return "", "", err
	}

	cipherData, err := symmetric.EncryptData([]byte(message), []byte(key))
	if err != nil {
		return "", "", err
	}

	cipherKey, err = EncryptMessageWithPublicKey(key)
	if err != nil {
		return
	}

	return base64.StdEncoding.EncodeToString(cipherData), cipherKey, nil
}

// DecryptLargeMessageWithPrivateKey decrypts the given base64 encoded message by
// decrypting the base64 encoded key with the rsa private key and then using
// the result to decrypt the ciphertext
func DecryptLargeMessageWithPrivateKey(message, key string) (string, error) {
	keyString, err := DecryptMessageWithPrivateKey(key)
	if err != nil {
		return "", fmt.Errorf("Unable to decrypt key with private key: %v", err)
	}

	messageData, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return "", err
	}

	data, err := symmetric.DecryptData(messageData, []byte(keyString))
	if err != nil {
		return "", err
	}

	return string(data), nil
}
