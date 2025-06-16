package Token

import (
	"log"
	TokenServices "v1/Services/Token"

	"github.com/golang-jwt/jwt/v5"
)

func (*HelperStruct) BlockToken(token *jwt.Token) error {
	if token == nil {
		return ErrNoToken
	}
	return TokenServices.SetTokenService(token)
}
func (*HelperStruct) IsBlocklisted(token *jwt.Token) bool {
	if token == nil {
		return false
	}
	tokenMap, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Println("Token cannot be mapped")
		return true
	}
	tokenId, ok := tokenMap["jti"].(string)
	if !ok {
		log.Default().Println("Token id is not present in token claims")
		return true
	}
	return TokenServices.GetTokenService(tokenId) != ""
}
