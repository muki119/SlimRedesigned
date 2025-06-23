package Controllers

import (
	"net/http"
	"v1/Helpers/Response"
	"v1/Helpers/Token"
	"v1/Middleware"

	"github.com/golang-jwt/jwt/v5"
)

func LogoutHandler(res http.ResponseWriter, req *http.Request) error {
	// if refresh token - add to blocked list -- service job
	// remove refresh token cookie
	parsedRefreshToken := req.Context().Value(Middleware.RequestTokenContextKey).(*jwt.Token)
	if parsedRefreshToken == nil {
		return ErrNoToken
	}
	err := Token.Token.BlockToken(parsedRefreshToken) // add parsed token id to blocked list
	if err != nil {
		return err
	}
	Response.SendCookieResponse(res, Response.ClearCookie(Response.RefreshTokenName), Response.SuccessResponse{Message: "successfully logged out"}, http.StatusOK)
	return nil
}
