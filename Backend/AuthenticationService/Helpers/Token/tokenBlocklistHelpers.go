package Token

import (
	TokenServices "v1/Services/Token"

	"github.com/golang-jwt/jwt/v5"
)

func (*HelperStruct) BlockToken(token *jwt.Token) error {
	return TokenServices.SetTokenService(token)
}
func (*HelperStruct) IsBlocklisted(token *jwt.Token) bool {
	tokenId := token.Claims.(jwt.MapClaims)["jti"].(string)
	return TokenServices.GetTokenService(tokenId) != ""
}
