package Services

import (
	"v1/Helpers/Password"
	"v1/Models"
)

func (userRepo *Services) RegisterService(user *Models.User) error { // get user data and create a new user // return only error response
	doesUserExist, err := user.UserExists()
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
	err = user.SaveUser()
	if err != nil {
		return err // return error if there is an issue saving the user
	}
	return nil
}
