package Helpers

import (
	"crypto/rsa"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JsonWebToken struct {
	Token string `json:"token"`
}

func CreateToken(userId, username string) (string, error) {
	// create a jwt thats signed with a private key
	audienceClaim := jwt.ClaimStrings{username}
	registeredClaims := jwt.RegisteredClaims{
		Subject:   userId,
		Audience:  audienceClaim,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 3)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "auth-service",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, registeredClaims)
	privateKey, err := getPrivateKey()
	if err != nil {
		return "", err
	}
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func getPrivateKey() (*rsa.PrivateKey, error) {
	// open pem file and load the private key
	//return []byte("your-private-key")
	publicKeyFileDir := os.Getenv("JWT_PUBLIC_KEY")
	fmt.Print("publicKeyFileDir: ", publicKeyFileDir)
	publicKeyFile, err := os.ReadFile(publicKeyFileDir) // pem file
	if err != nil {
		return nil, err
	}
	decodedPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM(publicKeyFile)
	if err != nil {
		return nil, err
	}
	return decodedPrivateKey, nil
}
