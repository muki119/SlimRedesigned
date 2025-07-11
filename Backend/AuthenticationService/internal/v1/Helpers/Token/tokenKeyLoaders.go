package Token

import (
	"crypto/rsa"
	"encoding/base64"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

// GetRSAPrivateKey Gets the Private key for the access token,
// Public key will only be used for the authentication by services other than the auth service
// Environment Key value must be a path to the .pem file
func GetRSAPrivateKey(privateKeyFileDir string) (*rsa.PrivateKey, error) {
	// open pem file and load the private envKey
	if privateKeyFileDir == "" {
		return nil, ErrNoPrivateKey
	}
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

// GetHMACSymmetricKey Gets the symmetric secret key for the Refresh token
// Environment Key value must be the
func GetHMACSymmetricKey(symmetricKey string) ([]byte, error) {
	if symmetricKey == "" {
		return nil, ErrNoSymmetricKey
	}
	symmetricKeyBytes, err := base64.StdEncoding.DecodeString(symmetricKey)
	if err != nil {
		return nil, err
	}
	return symmetricKeyBytes, nil
}
