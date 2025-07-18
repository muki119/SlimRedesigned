package Controllers

import (
	"AuthenticationService/v1/Helpers/Token"
	"AuthenticationService/v1/Services"
	"github.com/go-playground/validator/v10"
)

type Controllers struct {
	UserServices *Services.Services
	TokenHelpers *Token.Helper
	Validator    *validator.Validate
}
