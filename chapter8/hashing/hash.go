package hashing

import (
	"crypto/rand"
	"crypto/sha512"
	"fmt"
	"math/big"

	"github.com/dikaeinstein/go-microservice/chapter8/crypto"
)

// Hash is a structure which is capable of generating and comparing sha512
// hashes derived from a string
type Hash struct {
	peppers []string
}

// New creates a new hash and seeds with a list of peppers
func New(peppers []string) *Hash {
	return &Hash{peppers}
}

// GenerateHash hashes input string spiced up with salt and pepper using sha512.
func (h *Hash) GenerateHash(input string, withSalt, withPepper bool) (hash, salt string, err error) {
	var pepper string

	if withSalt {
		salt, err = GenerateRandomSalt()
		if err != nil {
			return
		}
	}

	if withPepper {
		pepper, err = h.getRandomPepper()
		if err != nil {
			return
		}
	}

	hash, err = h.createHash(input, salt, pepper)

	return
}

// Compare checks the input string against the salted hash. It'll loop through
// all peppers until finding success if the withPepper option is passed
func (h *Hash) Compare(input, salt, hash string, withPepper bool) (bool, error) {
	if withPepper {
		for _, p := range h.peppers {
			created, err := h.createHash(input, salt, p)
			if err != nil {
				return false, err
			}

			if created == hash {
				return true, nil
			}
		}
	} else {
		created, err := h.createHash(input, salt, "")
		if err != nil {
			return false, err
		}

		return created == hash, nil
	}

	return false, nil
}

func (h *Hash) getRandomPepper() (string, error) {
	max := big.NewInt(int64(len(h.peppers)))

	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}

	return h.peppers[n.Int64()], nil
}

func (h *Hash) createHash(input, salt, pepper string) (string, error) {
	stringToHash := input + salt + pepper

	sha := sha512.New()
	_, err := sha.Write([]byte(stringToHash))
	if err != nil {
		return "", err
	}

	hash := sha.Sum(nil)
	return fmt.Sprintf("%x", hash), nil
}

// GenerateRandomSalt generates random 32 bytes long salt
func GenerateRandomSalt() (string, error) {
	return crypto.GenerateRandomString(32)
}
