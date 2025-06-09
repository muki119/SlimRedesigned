package Token

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func (tokenService *HelperStruct) CreateAccessToken(userId string, issuer string) (string, error) { // access token - to be changed often -- rsa for microservices
	audienceClaim := jwt.ClaimStrings{"localhost:5000"}
	registeredClaims := jwt.RegisteredClaims{
		Subject:   userId,
		Audience:  audienceClaim,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    issuer,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, registeredClaims)
	privateKey := tokenService.PrivateKey
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
