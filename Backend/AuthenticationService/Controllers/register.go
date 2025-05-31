package Controllers

import (
	"fmt"
	"net/http"
)

type RegisterController struct {
	Forename      string
	Surname       string
	Username      string
	Password      string
	Date_of_birth string
}

func RegisterHandler(res http.ResponseWriter, req *http.Request) error {
	fmt.Println("RegisterHandler called")
	// Services.RegisterService()

	return nil
}
