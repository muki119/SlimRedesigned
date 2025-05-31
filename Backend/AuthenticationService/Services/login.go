package Services

type LoginServiceResponse struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

func LoginService() (*LoginServiceResponse, error) { // logs in a user and return their id and username for jwt creation
	// Get user credentials
	// Check if user exists
	// If yes, validate password
	// Return success or error response
	// return username and id
	return nil, nil
}
