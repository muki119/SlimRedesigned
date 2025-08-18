package token

import (
	"crypto/rsa"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type Token struct {
	// get the token asymmetrical key for jwt verificaion
	PublicKey *rsa.PublicKey
}

type Config struct {
	PublicKeyPath string
}

func (tokenConfig *Config) NewTokenHelper() (*Token, error) {
	publicKeyPemFile, err := os.ReadFile(tokenConfig.PublicKeyPath)
	if err != nil {
		return nil, err
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyPemFile)
	if err != nil {
		return nil, err
	}
	return &Token{
		PublicKey: publicKey,
	}, nil
}

func (tokenHelper *Token) VerifyToken(tokenString string) (bool, error) {
	token, err := tokenHelper.ParseToken(tokenString)
	if err != nil {
		return false, err
	}

	return token.Valid, nil // this verifies the token using the public key
}

func (tokenHelper *Token) ParseToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return tokenHelper.PublicKey, nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil // this parses the token and returns the token object
}
