package Services

import (
	"errors"
	"v1/Helpers/Password"
)

type LoginServiceResponse struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

func (userRepo *Services) LoginService(username string, password string) (*LoginServiceResponse, error) { // logs in a user and return their id and username for jwt creation
	userInfo, err := userRepo.UserRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	var validPassword bool
	validPassword, err = Password.ComparePassword(password, userInfo.Password)
	if err != nil {
		return nil, err
	}
	if !validPassword {
		return nil, ErrInvalidCredentials
	}
	return &LoginServiceResponse{
		Id:       userInfo.Id,
		Username: userInfo.Username,
	}, nil
}
