package jwt

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// CustomCliams represents the cliams used to create a JWT token
type CustomCliams struct {
	UserID      string `json:"userId"`
	AccessLevel string `json:"accessLevel"`
	jwt.StandardClaims
}

// GenerateJWT creates a new JWT and signs it with the private key
func GenerateJWT() []byte {
	rsaPrivate := initializePrivateKey()
	cliams := CustomCliams{
		"dikaeinstein",
		"user",
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			Issuer:    "dika",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, cliams)

	tokenString, err := token.SignedString(rsaPrivate)
	if err != nil {
		panic(err)
	}

	return []byte(tokenString)
}

// ValidateJWT validates that the given slice is a valid JWT and the signature matches
// the public key
func ValidateJWT(token []byte) error {
	rsaPublic := initializePublicKey()
	tokenString := string(token)

	parsedToken, err := jwt.ParseWithClaims(
		tokenString, &CustomCliams{},
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return rsaPublic, nil
		},
	)
	if err != nil {
		return fmt.Errorf("Unable to parse token: %v", err)
	}
	if !parsedToken.Valid {
		return errors.New("Unable to validate token")
	}

	return nil
}

func initializePrivateKey() *rsa.PrivateKey {
	var err error
	b, err := ioutil.ReadFile("../keys/sample_key.priv")
	if err != nil {
		log.Fatal("Unable to read private key file", err)
	}
	rsaPrivate, err := jwt.ParseRSAPrivateKeyFromPEM(b)
	if err != nil {
		log.Fatal("Unable to parse private key")
	}

	return rsaPrivate
}

func initializePublicKey() *rsa.PublicKey {
	var err error
	b, err := ioutil.ReadFile("../keys/sample_key.pub")
	if err != nil {
		log.Fatal("Unable to read public key file", err)
	}
	rsaPublic, err := jwt.ParseRSAPublicKeyFromPEM(b)
	if err != nil {
		log.Fatal("Unable to read public key file", err)
	}

	return rsaPublic
}
