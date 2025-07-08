package Services

import (
	"v1/Helpers/Password"
	"v1/dtos"
)

func (userRepo *Services) RegisterService(user *dtos.RegisterRequest) error { // get user data and create a new user // return only error response

	NewUser := userRepo.UserRepository.NewUser()
	NewUser.Forename = user.Forename
	NewUser.Surname = user.Surname
	NewUser.Username = user.Username
	NewUser.Email = user.Email
	NewUser.Password = user.Password
	NewUser.DateOfBirth = user.DateOfBirth

	doesUserExist, err := NewUser.UserExists()
	if doesUserExist {
		return err // if user exists - its going to return and error that contains what the matching field is , either email or username
	}
	if err != nil {
		return err
	}
	out, err := Password.HashPassword(user.Password)
	user.Password = out
	if err != nil {
		return err
	}
	err = NewUser.SaveUser()
	if err != nil {
		return err // return error if there is an issue saving the user
	}
	return nil
}
