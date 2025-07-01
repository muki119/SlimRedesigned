package Models

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"time"
	"v1/Config"
)

// hybrid approach centalised db for production, for testing use other function that takes in

// base user with names , username, email,password and role for the basic creations because dont need to set active , date created and last login

// repository gives out the new user and passes the database the repository will use for it

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
	db          *pgxpool.Pool
}

type ErrUserExists struct {
	Field   string `json:"field"`
	Message string `json:"error"`
}

var (
	ErrEmailExists    = &ErrUserExists{Field: "EMAIL", Message: "email already in use"}
	ErrUsernameExists = &ErrUserExists{Field: "USERNAME", Message: "username already in use"}
	ErrBothExists     = &ErrUserExists{Field: "BOTH", Message: "username and Email already in use"}
	ErrUserNotFound   = errors.New("user not found")
)

func (err ErrUserExists) Error() string {
	return err.Message
}

func (userRepo *UserRepository) CreateUserTable() error {
	_, err := userRepo.Db.Exec(Config.DatabaseContext, `
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
			role VARCHAR(30) NOT NULL DEFAULT 'USER' CHECK (UPPER(role) IN ('USER','ADMIN')),
			active BOOLEAN NOT NULL DEFAULT TRUE
		);
	`)
	if err != nil {
		return err
	}
	return nil
}

func (userRepo *UserRepository) NewUser() *User {
	return &User{
		Role: "USER",
		db:   userRepo.Db,
	}
}

func (u *User) SaveUser() error { // save user to database but if theres an issue roll back

	userDateOfBirth, err := time.Parse(time.RFC3339, u.DateOfBirth)
	if err != nil {
		return err
	}
	saveUserTx, err := u.db.Begin(Config.DatabaseContext)

	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			slog.Info("Rolling back save transaction", "error", p)
			err := saveUserTx.Rollback(Config.DatabaseContext)
			if err != nil {
				slog.Error("error rolling back.", "error", err.Error())
				return
			}
		}
		if err != nil {
			slog.Info("Rolling back save transaction", "error", err.Error())
			err := saveUserTx.Rollback(Config.DatabaseContext)
			if err != nil {
				slog.Error("error rolling back.", "error", err.Error())
				return
			}
		}
	}()
	_, err = saveUserTx.Exec(Config.DatabaseContext, `
		INSERT INTO users (forename, surname, username, email,password, role,date_of_birth)
		VALUES ($1, $2, $3, $4, $5, $6,$7)
	`, u.Forename, u.Surname, u.Username, u.Email, u.Password, u.Role, userDateOfBirth)
	if err != nil {
		return err
	}

	err = saveUserTx.Commit(Config.DatabaseContext)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) Delete() error {
	if u.Id == "" {
		return ErrUserNotFound
	}
	deleteUserTx, err := u.db.Begin(Config.DatabaseContext)
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			err := deleteUserTx.Rollback(Config.DatabaseContext)
			slog.Info("Rolling back delete transaction", "error", p)
			if err != nil {
				slog.Error("error rolling back.", "error", err.Error())
				return
			}
		}
		if err != nil {
			slog.Info("Rolling back delete transaction", "error", err.Error())
			err := deleteUserTx.Rollback(Config.DatabaseContext)
			if err != nil {
				return
			}
		}
	}()
	_, err = deleteUserTx.Exec(Config.DatabaseContext, `
		DELETE FROM users WHERE id = $1
	`, u.Id)

	err = deleteUserTx.Commit(Config.DatabaseContext)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) UserExists() (bool, error) {
	exists := u.db.QueryRow(Config.DatabaseContext, `
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
		return true, ErrBothExists
	} else if existingUser.Username == u.Username {
		return true, ErrUsernameExists
	} else if existingUser.Email == u.Email {
		return true, ErrEmailExists
	}
	return true, err
}

func (userRepo *UserRepository) GetUserByUsername(username string) (*User, error) {
	// Implement the logic to retrieve a user by username from the database
	// This is a placeholder function and should be replaced with actual database logic
	userData := userRepo.Db.QueryRow(context.Background(), `
		SELECT id, forename, surname, username, email, password,date_of_birth::text, date_created::text, last_login::text, role, active From users WHERE username=$1
	`, username)

	var foundUserData = userRepo.NewUser()
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
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return foundUserData, err
}

func (userRepo *UserRepository) GetUserByEmail(email string) (*User, error) {
	userData := userRepo.Db.QueryRow(context.Background(), `
		SELECT id, forename, surname, username, email, password, date_created, last_login, role, active From users WHERE email=$1
	`, email)

	var foundUserData = userRepo.NewUser()
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
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return foundUserData, err
}
