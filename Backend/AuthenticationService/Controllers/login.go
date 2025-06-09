package Controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"v1/Helpers/Response"
	"v1/Helpers/Token"
	"v1/Models"
	"v1/Services"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginHandler(res http.ResponseWriter, req *http.Request) error {
	var loginReq LoginRequest
	json.NewDecoder(req.Body).Decode(&loginReq)
	var user, err = Services.LoginService(loginReq.Username, loginReq.Password) // should return the id and
	if err != nil {
		if errors.Is(err, Models.UserNotFoundError) {
			Response.SendJsonResponse(res, &Response.ErrorResponse{Error: err.Error()}, http.StatusNotFound)
			return nil
		}
		return err
	}
	// make refresh token
	refreshToken, err := Token.Token.CreateLoginRefreshToken(user.Id)
	if err != nil {
		return err
	}
	accessToken, err := Token.Token.CreateAccessToken(user.Id, "/login")
	if err != nil {
		return err
	}

	Response.SendCookieResponse(res, Response.NewRefreshTokenCookie(refreshToken), &Response.AccessTokenResponse{Token: accessToken, Message: "successfully logged in"}, http.StatusOK)
	return nil
}
