package Controllers

import (
	"errors"
	"net/http"
	"v1/Helpers/Response"
	"v1/Helpers/Token"
)

func TokenHandler(res http.ResponseWriter, req *http.Request) error {
	// get a refresh token
	requestRefreshToken, err := req.Cookie("RefreshToken")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			Response.SendJsonResponse(res, Response.ErrorResponse{Error: "Token Invalid"}, http.StatusUnauthorized)
			return nil
		}
		return err
	}
	parsedRefreshToken, err := Token.Token.ParseToken(requestRefreshToken.Value) // check if the jwt is valid
	if err != nil {
		return err
	}
	err = Token.Token.BlockToken(parsedRefreshToken) // add parsed token id to blocked list
	if err != nil {
		return err
	}
	newRefreshToken, err := Token.Token.CreateRefreshTokenFromClaims(parsedRefreshToken.Claims)
	if err != nil {
		return err
	}
	userId, err := parsedRefreshToken.Claims.GetSubject()
	if err != nil {
		return err
	}
	newAccessToken, err := Token.Token.CreateAccessToken(userId, "/token")
	if err != nil {
		return err
	}

	Response.SendCookieResponse(res, Response.NewRefreshTokenCookie(newRefreshToken), Response.AccessTokenResponse{
		Message: "successfully refreshed",
		Token:   newAccessToken,
	}, http.StatusCreated)

	return nil
}
