package Controllers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"v1/Helpers/Response"
	"v1/Models"
)

func (controllerTools *Controllers) RegisterHandler(res http.ResponseWriter, req *http.Request) error {
	userDetails := controllerTools.UserServices.UserRepository.NewUser() // stores the user's details to be passed
	err := json.NewDecoder(req.Body).Decode(&userDetails)
	if err != nil { // if there's an error decoding
		log.Println(err)
		return err
	}
	err = controllerTools.UserServices.RegisterService(userDetails) // will return an error if the user already exists or if there is an issue saving the user
	if err != nil {
		var ErrUserExists *Models.ErrUserExists // if theres an error regestring the user
		if errors.As(err, &ErrUserExists) {     // and the error comes from
			Response.SendJsonResponse(res, err, http.StatusBadRequest)
			return nil
		}
		return err
	}
	Response.SendJsonResponse(res, &Response.SuccessResponse{Message: "Successfully registered"}, http.StatusCreated)
	return nil
}
