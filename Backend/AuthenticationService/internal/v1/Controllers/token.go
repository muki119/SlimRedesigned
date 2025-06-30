package Controllers

import (
	"errors"
	"log"
	"net/http"
	"v1/Helpers/Response"
	"v1/Helpers/Token"
	"v1/Middleware"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrNoToken = errors.New("no token found")
)

func (*Controllers) TokenHandler(res http.ResponseWriter, req *http.Request) error {
	// get a refresh token
	parsedRefreshToken := req.Context().Value(Middleware.RequestTokenContextKey).(*jwt.Token)
	if parsedRefreshToken == nil {
		return ErrNoToken
	}
	newRefreshToken, err := Token.Token.CreateRefreshTokenFromClaims(parsedRefreshToken.Claims)
	if err != nil {
		log.Println("error in Token.Token.CreateRefreshTokenFromClaims", err)
		return err
	}
	userId, err := parsedRefreshToken.Claims.GetSubject()
	if err != nil {
		log.Println("error parsing refreshtoken", err)
		return err
	}
	newAccessToken, err := Token.Token.CreateAccessToken(userId, "/token")
	if err != nil {
		log.Println("error in Token.CreateAccessToken: ", err)
		return err
	}

	err = Token.Token.BlockToken(parsedRefreshToken) // add parsed token id to blocked list
	if err != nil {
		return err
	}
	Response.SendCookieResponse(res, Response.NewRefreshTokenCookie(newRefreshToken), Response.AccessTokenResponse{
		Message: "successfully refreshed",
		Token:   newAccessToken,
	}, http.StatusOK)

	return nil
}
