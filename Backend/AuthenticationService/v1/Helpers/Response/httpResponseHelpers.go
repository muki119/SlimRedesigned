package Response

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

type ErrorResponse struct {
	Error any `json:"error"`
}
type SuccessResponse struct {
	Message string `json:"message"`
}
type AccessTokenResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

const RefreshTokenName = "refresh_token"

func NewRefreshTokenCookie(value string) *http.Cookie {
	return &http.Cookie{
		Name:     RefreshTokenName,
		HttpOnly: true,
		Path:     "/",
		Value:    value,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
}

func ClearCookie(cookieName string) *http.Cookie {
	return &http.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	}
}

func SendJsonResponse(res http.ResponseWriter, outStruct any, statusCode int) {
	res.Header().Set("Server", "H26")
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.WriteHeader(statusCode)
	err := json.NewEncoder(res).Encode(outStruct)
	if err != nil {
		slog.Error(err.Error())
		return
	}
}

func SendCookieResponse(res http.ResponseWriter, cookie *http.Cookie, outStruct any, statusCode int) {
	http.SetCookie(res, cookie)
	res.Header().Set("Server", "H26")
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.WriteHeader(statusCode)
	err := json.NewEncoder(res).Encode(outStruct)
	if err != nil {
		slog.Error(err.Error())
		return
	}
}
