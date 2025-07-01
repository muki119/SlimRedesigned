package Controllers

import (
	"v1/Helpers/Token"
	"v1/Services"
)

type Controllers struct {
	UserServices *Services.Services
	TokenHelpers *Token.Helper
}
