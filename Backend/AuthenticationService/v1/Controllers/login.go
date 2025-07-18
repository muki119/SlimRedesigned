package Controllers

import (
	"AuthenticationService/v1/Helpers/Response"
	"AuthenticationService/v1/Models"
	"AuthenticationService/v1/Services"
	"AuthenticationService/v1/Utils"
	"AuthenticationService/v1/dtos"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func (controllerTools *Controllers) LoginHandler(res http.ResponseWriter, req *http.Request) error {
	var loginReq dtos.LoginRequest
	err := json.NewDecoder(req.Body).Decode(&loginReq)
	if err != nil {
		return err
	}

	if err = controllerTools.Validator.Struct(loginReq); err != nil {
		err := err.(validator.ValidationErrors)
		Response.SendJsonResponse(res, &Response.ErrorResponse{Error: Utils.FormatErrors(err)}, http.StatusUnprocessableEntity)
		return nil
	}
	user, err := controllerTools.UserServices.LoginService(loginReq.Username, loginReq.Password) // should return the id and
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
	refreshToken, err := controllerTools.TokenHelpers.CreateLoginRefreshToken(user.Id)
	if err != nil {
		return err
	}
	//make access token
	// the access token will be used to access protected routes
	accessToken, err := controllerTools.TokenHelpers.CreateAccessToken(user.Id, "/login")
	if err != nil {
		return err
	}

	Response.SendCookieResponse(res, Response.NewRefreshTokenCookie(refreshToken), &Response.AccessTokenResponse{Token: accessToken, Message: "successfully logged in"}, http.StatusOK)
	return nil
}
