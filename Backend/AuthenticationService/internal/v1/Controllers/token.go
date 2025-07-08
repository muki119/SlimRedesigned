package Controllers

import (
	"errors"
	"log/slog"
	"net/http"
	"v1/Helpers/Response"
	"v1/Middleware"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrNoToken = errors.New("no token found")
)

func (controllerTools *Controllers) TokenHandler(res http.ResponseWriter, req *http.Request) error {
	// get a refresh token
	parsedRefreshToken := req.Context().Value(Middleware.RequestTokenContextKey).(*jwt.Token)
	if parsedRefreshToken == nil {
		return ErrNoToken
	}
	newRefreshToken, err := controllerTools.TokenHelpers.CreateRefreshTokenFromClaims(parsedRefreshToken.Claims)
	if err != nil {
		slog.Error("error in Token.Token.CreateRefreshTokenFromClaims", "error", err.Error())
		return err
	}
	userId, err := parsedRefreshToken.Claims.GetSubject()
	if err != nil {
		slog.Error("error parsing refresh token", "error", err.Error())
		return err
	}
	newAccessToken, err := controllerTools.TokenHelpers.CreateAccessToken(userId, "/token")
	if err != nil {
		slog.Error("error in Token.CreateAccessToken: ", "error", err.Error())
		return err
	}

	err = controllerTools.TokenHelpers.Blocklist.BlockToken(parsedRefreshToken) // add parsed token id to blocked list
	if err != nil {
		return err
	}
	Response.SendCookieResponse(res, Response.NewRefreshTokenCookie(newRefreshToken), Response.AccessTokenResponse{
		Message: "successfully refreshed",
		Token:   newAccessToken,
	}, http.StatusOK)

	return nil
}
