package Controllers

import (
	"fmt"
	"net/http"
)

func LogoutHandler(res http.ResponseWriter, req *http.Request) error {
	fmt.Println("LogoutHandler called")
	res.Write([]byte("{error:'yoo'}"))
	//replace their jwt token - which should be in a cookie form with nothing.
	// Services.LogoutService()

	return nil
}
