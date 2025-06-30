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

func (controllers *Controllers) LoginHandler(res http.ResponseWriter, req *http.Request) error {
	var loginReq LoginRequest
	err := json.NewDecoder(req.Body).Decode(&loginReq)
	if err != nil {
		return err
	}
	user, err := controllers.UserServices.LoginService(loginReq.Username, loginReq.Password) // should return the id and
	if err != nil {
		if errors.Is(err, Models.ErrUserNotFound) {
			Response.SendJsonResponse(res, &Response.ErrorResponse{Error: err.Error()}, http.StatusNotFound)
			return nil
		}
		if errors.Is(err, Services.ErrInvalidCredentials) {
			Response.SendJsonResponse(res, &Response.ErrorResponse{Error: err.Error()}, http.StatusUnauthorized)
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
