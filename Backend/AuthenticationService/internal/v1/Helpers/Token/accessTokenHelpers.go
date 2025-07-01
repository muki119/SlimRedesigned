package Token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (tokenService *Helper) CreateAccessToken(userId string, issuer string) (string, error) { // access token - to be changed often -- rsa for microservices
	if userId == "" {
		return "", ErrNoUserId
	} else if issuer == "" {
		return "", ErrNoIssuer
	}

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
