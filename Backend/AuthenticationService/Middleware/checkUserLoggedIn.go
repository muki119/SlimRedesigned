package Middleware

import (
	"errors"
	"net/http"
	"v1/Helpers/Response"
	"v1/Helpers/Token"
)

func CheckUserLoggedIn(f func(http.ResponseWriter, *http.Request) error) func(http.ResponseWriter, *http.Request) error {
	return func(res http.ResponseWriter, req *http.Request) error {
		//check if user is logged in by verifying jwt signature is verified and has a subject
		// verify jwt and that its not blocklisted
		requestRefreshToken, err := req.Cookie("RefreshToken")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) { // if no cookie
				Response.SendJsonResponse(res, Response.ErrorResponse{Error: "Token Invalid"}, http.StatusUnauthorized) // send invalid cookie
				return nil
			}
			return err
		}
		parsedRefreshToken, err := Token.Token.ParseToken(requestRefreshToken.Value)
		if err != nil {
			return err
		}
		if Token.Token.IsBlocklisted(parsedRefreshToken) || !parsedRefreshToken.Valid { // check blocklisted and not valid
			Response.SendCookieResponse(res, Response.ClearCookie("RefreshToken"), Response.ErrorResponse{Error: "Token Invalid"}, http.StatusUnauthorized)
			return nil
		}
		err = f(res, req)
		return err
	}
}
