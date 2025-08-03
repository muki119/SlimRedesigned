package Utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type FieldError struct {
	Tag string `json:"tag"`
	Msg string `json:"msg"`
}

func (fe FieldError) Error() string {
	return fe.Msg
}

func FormatValidationErrors(errs validator.ValidationErrors) *map[string][]*FieldError {
	errorMap := make(map[string][]*FieldError)
	for _, fieldError := range errs {
		errorMap[fieldError.Field()] = append(errorMap[fieldError.Field()], generateValidationErrors(fieldError))
	}
	return &errorMap
}

func generateValidationErrors(fieldError validator.FieldError) *FieldError {
	switch fieldError.Tag() {
	case "required":
		return &FieldError{
			Tag: fieldError.Tag(),
			Msg: fmt.Sprintf("%s is required", fieldError.Field())}
	case "min":
		return &FieldError{
			Tag: fieldError.Tag(),
			Msg: fmt.Sprintf("%s needs to be at least %s ", fieldError.Field(), fieldError.Param()),
		}
	case "max":
		return &FieldError{
			Tag: fieldError.Tag(),
			Msg: fmt.Sprintf("%s needs to be at most %s ", fieldError.Field(), fieldError.Param()),
		}
	case "email":
		return &FieldError{
			Tag: fieldError.Tag(),
			Msg: "value is not a valid email address",
		}
	case "alpha":
		return &FieldError{
			Tag: fieldError.Tag(),
			Msg: "value must contain only alphabetic characters",
		}
	case "alphanum":
		return &FieldError{
			Tag: fieldError.Tag(),
			Msg: "value must contain only alphanumeric characters",
		}
	default:
		return &FieldError{
			Tag: fieldError.Tag(),
			Msg: fmt.Sprintf("%s is invalid", fieldError.Field()),
		}
	}
}
