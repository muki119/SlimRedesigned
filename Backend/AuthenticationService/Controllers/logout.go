package Controllers

import (
	"errors"
	"net/http"
	"v1/Helpers/Response"
	"v1/Helpers/Token"
)

func LogoutHandler(res http.ResponseWriter, req *http.Request) error {
	// if refresh token - add to blocked list -- service job
	// remove refresh token cookie
	requestRefreshToken, err := req.Cookie("RefreshToken")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) { // if no cookie
			Response.SendJsonResponse(res, Response.ErrorResponse{Error: "Token Invalid"}, http.StatusUnauthorized) // send invalid cookie
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
	Response.SendCookieResponse(res, Response.ClearCookie("RefreshToken"), Response.ErrorResponse{Error: "successfully logged out"}, http.StatusNoContent)
	return nil
}
