package Controllers

import (
	"AuthenticationService/v1/Helpers/Response"
	"AuthenticationService/v1/Models"
	"AuthenticationService/v1/Utils"
	"AuthenticationService/v1/dtos"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func (controllerTools *Controllers) RegisterHandler(res http.ResponseWriter, req *http.Request) error {

	var userDetails *dtos.RegisterRequest // stores the user's details to be passed
	err := json.NewDecoder(req.Body).Decode(&userDetails)
	if err != nil { // if there's an error decoding
		return err
	}
	if err = controllerTools.Validator.Struct(userDetails); err != nil {
		err := err.(validator.ValidationErrors)
		Response.SendJsonResponse(res, &Response.ErrorResponse{Error: Utils.FormatErrors(err)}, http.StatusUnprocessableEntity)
		return nil
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
