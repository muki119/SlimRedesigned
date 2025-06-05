package Helpers

import (
	"encoding/json"
	"net/http"
)

func SendJsonResponse(res http.ResponseWriter, outStruct any, statusCode int) {
	res.Header().Set("Server", "H26")
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.WriteHeader(statusCode)
	json.NewEncoder(res).Encode(outStruct)
}

func SendCookieResponse(res http.ResponseWriter, cookie *http.Cookie, outStruct any, statusCode int) {
	http.SetCookie(res, cookie)
	res.Header().Set("Server", "H26")
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.WriteHeader(statusCode)
	json.NewEncoder(res).Encode(outStruct)
}
