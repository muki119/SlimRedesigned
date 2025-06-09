package Token

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

func (tokenService *HelperStruct) createRefreshToken(userId string, expiresAt *jwt.NumericDate, issuer string) (string, error) { // for auth server use only
	audienceClaim := jwt.ClaimStrings{"localhost:5000"}
	jti := uuid.New().String()
	if jti == "" {
		return "", fmt.Errorf("jwt Id could not be generated")
	}
	registeredClaims := jwt.RegisteredClaims{
		Subject:   userId,
		Audience:  audienceClaim,
		ExpiresAt: expiresAt,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    issuer,
		ID:        jti,
	}
	secretKey := tokenService.SymmetricKey
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, registeredClaims)
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
func (tokenService *HelperStruct) CreateLoginRefreshToken(userId string) (string, error) {
	expiresAt := jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7))
	return tokenService.createRefreshToken(userId, expiresAt, "/login")
}
func (tokenService *HelperStruct) CreateRefreshTokenFromClaims(tokenClaims jwt.Claims) (string, error) {
	userId, err := tokenClaims.GetSubject()
	if err != nil {
		return "", err
	}
	expiresAt, err := tokenClaims.GetExpirationTime()
	if err != nil {
		return "", err
	}
	return tokenService.createRefreshToken(userId, expiresAt, "/refresh")
}
