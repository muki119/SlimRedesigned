package Controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"v1/Helpers"
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
		if errors.Is(err, Models.UserNotFound) {
			Helpers.SendJsonResponse(res, &Helpers.ErrorResponse{Error: err.Error()}, http.StatusNotFound)
			return nil
		}
		return err
	}
	// make a jwt
	userToken, err := Helpers.CreateToken(user.Id, user.Username)
	if err != nil {
		return err
	}
	cookie := &http.Cookie{
		Name:     "SID",
		Value:    userToken,
		HttpOnly: true,
		Path:     "/",
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
	Helpers.SendCookieResponse(res, cookie, &Helpers.SuccessResponse{Message: "Successfully logged in"}, http.StatusOK)
	return nil
}
