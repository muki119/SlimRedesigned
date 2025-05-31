package Models

import (
	"context"
	"fmt"
	"v1/Config"
)

type User struct {
	Forename     string `json:"forename"`
	Surname      string `json:"surname"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Email        string `json:"email"`
	Date_Created string `json:"date_created"`
	Last_Login   string `json:"last_login"`
	Role         string `json:"role"`   // e.g., "admin", "user", etc.
	Active       bool   `json:"active"` // Indicates if the user account is active
}

type UserWithId struct {
	Id string `json:"id"` // Unique identifier for the user
	User
}

func CreateUserTable() error {
	out, err := Config.DatabaseConnection.Prepare(context.Background(), "CreateUserTable", `
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			forename VARCHAR(60) NOT NULL,
			surname VARCHAR(30) NOT NULL,
			username VARCHAR(30) NOT NULL UNIQUE,
			email VARCHAR(320) NOT NULL UNIQUE,
			password VARCHAR(100) NOT NULL,
			date_created TIMESTAMP NOT NULL DEFAULT NOW(),
			last_login TIMESTAMP NOT NULL DEFAULT NOW(),
			role VARCHAR(20) NOT NULL DEFAULT 'user',
			active BOOLEAN NOT NULL DEFAULT TRUE
		)
	`)
	fmt.Println(out)
	return err
}

func (u *User) NewUser(forename, surname, username, password, email, dateCreated, lastLogin, role string, active bool) User {
	return User{
		Forename:     forename,
		Surname:      surname,
		Username:     username,
		Password:     password,
		Email:        email,
		Date_Created: dateCreated,
		Last_Login:   lastLogin,
		Role:         role,
		Active:       active,
	}
}
func (u *User) SaveUser() error {
	// Implement the logic to save the user to the database
	// This is a placeholder function and should be replaced with actual database logic
	return nil
}

func GetUser(username string) (*UserWithId, error) {
	// Implement the logic to retrieve a user by username from the database
	// This is a placeholder function and should be replaced with actual database logic
	return nil, nil
}
