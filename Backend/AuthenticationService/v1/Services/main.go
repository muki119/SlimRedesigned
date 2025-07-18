package Services

import "AuthenticationService/v1/Models"

type Services struct {
	UserRepository *Models.UserRepository
}

type servicesInterface interface {
	RegisterService(user *Models.User) error
	LoginService(string, string) (*LoginServiceResponse, error)
}
