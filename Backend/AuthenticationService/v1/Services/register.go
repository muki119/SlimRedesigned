package Services

import (
	"AuthenticationService/v1/Helpers/Password"
	"AuthenticationService/v1/dtos"
)

func (userRepo *Services) RegisterService(registerRequest *dtos.RegisterRequest) error { // get user data and create a new user // return only error response

	NewUser := userRepo.UserRepository.NewUser()

	NewUser.Forename = registerRequest.Forename
	NewUser.Surname = registerRequest.Surname
	NewUser.Username = registerRequest.Username
	NewUser.Email = registerRequest.Email
	NewUser.Password = registerRequest.Password
	NewUser.DateOfBirth = registerRequest.DateOfBirth

	doesUserExist, err := NewUser.UserExists()
	if doesUserExist {
		return err // if user exists - its going to return and error that contains what the matching field is , either email or username
	}
	if err != nil {
		return err
	}
	out, err := Password.HashPassword(NewUser.Password)
	if err != nil {
		return err
	}
	NewUser.Password = out

	err = NewUser.SaveUser()
	if err != nil {
		return err // return error if there is an issue saving the user
	}
	return nil
}
