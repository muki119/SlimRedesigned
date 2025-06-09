package Token

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

func (tokenService *HelperStruct) ParseToken(token string) (*jwt.Token, error) {
	secretKey := tokenService.SymmetricKey
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})
}
