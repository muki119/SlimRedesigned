package Middleware

import (
	"net/http"
	"v1/Helpers/Token"
)

type Middleware struct {
	TokenHelper *Token.Helper
}

type middlewareInterface interface {
	CheckUserLoggedIn(func(http.ResponseWriter, *http.Request) error) func(http.ResponseWriter, *http.Request) error
	ErrorHandler(func(http.ResponseWriter, *http.Request) error) http.HandlerFunc
}
