package Token

import (
	"crypto/rsa"
	"encoding/base64"
	"github.com/golang-jwt/jwt/v5"
	"os"
)

func getRSAPrivateKey() (*rsa.PrivateKey, error) {
	// open pem file and load the private key
	privateKeyFileDir := os.Getenv("JWT_PRIVATE_KEY")
	privateKeyFile, err := os.ReadFile(privateKeyFileDir) // pem file
	if err != nil {
		return nil, err
	}
	decodedPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyFile)
	if err != nil {
		return nil, err
	}
	return decodedPrivateKey, nil
}

func getHMACSymmetricKey() ([]byte, error) {
	symmetricKey := os.Getenv("JWT_SECRET_KEY")
	symmetricKeyBytes, err := base64.StdEncoding.DecodeString(symmetricKey)
	if err != nil {
		return nil, err
	}
	return symmetricKeyBytes, nil
}
