package Controllers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"v1/Helpers/Response"
	"v1/Models"
	"v1/Services"
)

func RegisterHandler(res http.ResponseWriter, req *http.Request) error {
	userDetails := Models.NewUser() // stores the user's details to be passed
	err := json.NewDecoder(req.Body).Decode(&userDetails)
	if err != nil { // if there's an error decoding
		log.Println(err)
		return err
	}
	err = Services.RegisterService(userDetails) // will return an error if the user already exists or if there is an issue saving the user
	if err != nil {                             // if theres an error regestring the user
		if errors.As(err, &Models.UserExistsErrorPtr) { // and the error comes from
			var errorOut *Models.UserExistsError
			errors.As(err, &errorOut)
			Response.SendJsonResponse(res, errorOut, http.StatusBadRequest)
			return nil
		}
		return err
	}
	Response.SendJsonResponse(res, &Response.SuccessResponse{Message: "Successfully registered"}, http.StatusCreated)
	return nil
}
