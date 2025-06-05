package Services

import (
	"v1/Helpers"
	"v1/Models"
)

type LoginServiceResponse struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

func LoginService(username string, password string) (*LoginServiceResponse, error) { // logs in a user and return their id and username for jwt creation
	userInfo, err := Models.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	var validPassword bool
	validPassword, err = Helpers.ComparePassword(password, userInfo.Password)
	if !validPassword {
		return nil, err
	}
	return &LoginServiceResponse{
		Id:       userInfo.Id,
		Username: userInfo.Username,
	}, nil
}
