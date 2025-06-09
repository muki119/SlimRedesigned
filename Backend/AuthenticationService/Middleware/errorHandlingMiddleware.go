package Middleware

import (
	"log"
	"net/http"
	"v1/Helpers/Response"
)

func ErrorHandler(fn func(http.ResponseWriter, *http.Request) error) http.HandlerFunc { // a middleware that takes a httpp handler that also returns an error and handles it by send ing a 500 internal server error response
	return func(res http.ResponseWriter, req *http.Request) {
		err := fn(res, req)
		if err != nil {
			// send a 500 status with a error json response
			log.Println(err)
			Response.SendJsonResponse(res, &Response.ErrorResponse{Error: "internal server error"}, http.StatusInternalServerError)
		}
	}
}
