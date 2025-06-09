package Models

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"v1/Config"
)

//base user with names , username, email,password and role for the basic creations because dont need to set active , date created and last login

type User struct {
	Id          string `json:"id" db:"id"`
	Forename    string `json:"forename" db:"forename"`
	Surname     string `json:"surname"  db:"surname"`
	Username    string `json:"username"   db:"username"`
	Email       string `json:"email" db:"email"`
	Password    string `json:"password"  db:"password"`
	DateOfBirth string `json:"date_of_birth" db:"date_of_birth"`
	Role        string `json:"role"   db:"role"`
	DateCreated string `json:"date_created" db:"date_created"`
	LastLogin   string `json:"last_login" db:"last_login"`
	Active      bool   `json:"active"  db:"active"`
}

type UserExistsError struct {
	Field   string `json:"field"`
	Message string `json:"error"`
}

var (
	EmailExistsError    = &UserExistsError{Field: "EMAIL", Message: "email already in use"}
	UsernameExistsError = &UserExistsError{Field: "USERNAME", Message: "username already in use"}
	BothExistsError     = &UserExistsError{Field: "BOTH", Message: "username and Email already in use"}
	UserNotFoundError   = errors.New("user not found")
)

var UserExistsErrorPtr *UserExistsError

func (err UserExistsError) Error() string {
	return err.Message
}

func CreateUserTable() {
	//_, err := Config.DatabaseConnection.Exec(context.Background(), `
	//	DROP TYPE IF EXISTS USER_ROLE_TYPE;
	//	CREATE TYPE USER_ROLE_TYPE AS ENUM ('USER','ADMIN');
	//`)
	//if err != nil {
	//	panic(err)
	//}
	_, err := Config.DatabaseConnection.Prepare(Config.DatabaseContext, "CreateUserTable", `
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			forename VARCHAR(60) NOT NULL,
			surname VARCHAR(30) NOT NULL,
			username VARCHAR(30) NOT NULL UNIQUE,
			email VARCHAR(320) NOT NULL UNIQUE,
			password VARCHAR(100) NOT NULL,
			date_of_birth DATE NOT NULL,
			date_created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			last_login TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			role USER_ROLE_TYPE NOT NULL DEFAULT 'USER',
			active BOOLEAN NOT NULL DEFAULT TRUE
		);
	`)
	if err != nil {
		panic(err)
	}
}
func NewUser() *User {
	return &User{
		Role: "USER",
	}
}

func NewAdminUser() *User {
	return &User{
		Role: "ADMIN",
	}
}

func (u *User) SaveUser() error {
	saveUserTransaction, err := Config.DatabaseConnection.Begin(Config.DatabaseContext)
	if err != nil {
		return err
	}
	dobToDate, _ := time.Parse(time.DateOnly, u.DateOfBirth)
	_, err = saveUserTransaction.Exec(Config.DatabaseContext, `
		INSERT INTO users (forename, surname, username, email,password, role,date_of_birth)
		VALUES ($1, $2, $3, $4, $5, $6,$7)
	`, u.Forename, u.Surname, u.Username, u.Email, u.Password, u.Role, dobToDate)
	if err != nil {
		saveUserTransaction.Rollback(Config.DatabaseContext)
		return err
	}
	saveUserTransaction.Commit(Config.DatabaseContext)
	return nil
}

func (u *User) UserExists() (bool, error) {
	exists := Config.DatabaseConnection.QueryRow(Config.DatabaseContext, `
		SELECT username,email FROM users WHERE username=$1 OR email=$2;
	`, u.Username, u.Email)

	type output struct {
		Username string `db:"username"`
		Email    string `db:"email"`
	}
	var existingUser output

	err := exists.Scan(&existingUser.Username, &existingUser.Email)
	if err != nil { // if there was an error
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		} // and the error was because the username/ email wasnt in the database
		return false, err
	}
	if existingUser.Username == u.Username && existingUser.Email == u.Email {
		return true, BothExistsError
	} else if existingUser.Username == u.Username {
		return true, UsernameExistsError
	} else if existingUser.Email == u.Email {
		return true, EmailExistsError
	}
	return true, err
}

func GetUserByUsername(username string) (*User, error) {
	// Implement the logic to retrieve a user by username from the database
	// This is a placeholder function and should be replaced with actual database logic
	userData := Config.DatabaseConnection.QueryRow(context.Background(), `
		SELECT id, forename, surname, username, email, password,date_of_birth::text, date_created::text, last_login::text, role, active From users WHERE username=$1
	`, username)

	var foundUserData User
	err := userData.Scan(&foundUserData.Id,
		&foundUserData.Forename,
		&foundUserData.Surname,
		&foundUserData.Username,
		&foundUserData.Email,
		&foundUserData.Password,
		&foundUserData.DateOfBirth,
		&foundUserData.DateCreated,
		&foundUserData.LastLogin,
		&foundUserData.Role,
		&foundUserData.Active)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, UserNotFoundError
		}
		return nil, err
	}
	return &foundUserData, err
}

func GetUserByEmail(email string) (*User, error) {
	userData := Config.DatabaseConnection.QueryRow(context.Background(), `
		SELECT id, forename, surname, username, email, password, date_created, last_login, role, active From users WHERE email=$1
	`, email)

	var foundUserData User
	err := userData.Scan(&foundUserData.Id,
		&foundUserData.Forename,
		&foundUserData.Surname,
		&foundUserData.Username,
		&foundUserData.Email,
		&foundUserData.Password,
		&foundUserData.DateCreated,
		&foundUserData.LastLogin,
		&foundUserData.Role,
		&foundUserData.Active)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, UserNotFoundError
		}
		return nil, err
	}
	return &foundUserData, err
}
