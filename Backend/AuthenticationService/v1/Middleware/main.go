package Middleware

import (
	"AuthenticationService/v1/Helpers/Token"
	"net/http"
)

type Middleware struct {
	TokenHelper *Token.Helper
}

type middlewareInterface interface {
	CheckUserLoggedIn(func(http.ResponseWriter, *http.Request) error) func(http.ResponseWriter, *http.Request) error
	ErrorHandler(func(http.ResponseWriter, *http.Request) error) http.HandlerFunc
}
