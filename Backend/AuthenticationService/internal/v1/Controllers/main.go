package Controllers

import (
	"github.com/go-playground/validator/v10"
	"v1/Helpers/Token"
	"v1/Services"
)

type Controllers struct {
	UserServices *Services.Services
	TokenHelpers *Token.Helper
	Validator    *validator.Validate
}
