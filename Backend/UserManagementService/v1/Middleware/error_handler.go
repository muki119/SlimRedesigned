package middleware

import (
	"log/slog"
	"net/http"
	"v1/Utils"
)

func (*Middleware) ErrorHandler(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if err := f(res, req); err != nil {
			slog.Error(err.Error())
			Utils.SendJsonResponse(res, &Utils.ErrorResponse{Error: "Internal Server Error"}, http.StatusInternalServerError)
		}
	}
}
