package middleware

import (
	"context"
	"net/http"
	"strings"
	"v1/Utils"
)

// Takes a handler function that returns error and checks if the incomming request jwt is valid , then it initiates the passed function
type contextKey string

func (middlewareTools *Middleware) UserAuthMiddleware(next func(http.ResponseWriter, *http.Request) error) func(http.ResponseWriter, *http.Request) error { // this entire function is wrapped in a error handler
	return func(res http.ResponseWriter, req *http.Request) error {
		bearerToken := req.Header.Get("Authorization")
		tokenString := strings.Split(bearerToken, "Bearer ")
		if len(tokenString) != 2 {
			Utils.SendJsonResponse(res, Utils.ErrorResponse{Error: "Invalid Authorization contents"}, http.StatusBadRequest)
			return nil
		}
		parsedToken, err := middlewareTools.TokenHelper.ParseToken(tokenString[1])
		if err != nil {
			return err
		}
		if !parsedToken.Valid {
			Utils.SendJsonResponse(res, Utils.ErrorResponse{Error: "Invalid Authorization token"}, http.StatusBadRequest)
			return nil
		}
		newRequestContext := context.WithValue(req.Context(), contextKey("user"), parsedToken)

		return next(res, req.WithContext(newRequestContext))
	}
}
