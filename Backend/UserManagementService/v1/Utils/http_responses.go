package Utils

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type ErrorResponse struct {
	Error any `json:"error"`
}
type SuccessResponse struct {
	Message string `json:"message"`
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
