package Controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"v1/Services"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginHandler(res http.ResponseWriter, req *http.Request) error {
	var loginReq LoginRequest
	json.NewDecoder(req.Body).Decode(&loginReq)
	_, err := Services.LoginService()
	if err != nil {
		return err
	}
	// create a jwt with the rsa private key - only real thing that should be in there is the users Id
	// send the jwt only

	fmt.Printf("Username: %s, Password: %s\n", loginReq.Username, loginReq.Password)
	// Services.LoginService()
	return nil
}
