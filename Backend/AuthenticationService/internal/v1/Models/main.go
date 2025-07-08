package Models

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

type UserRepository struct {
	Db *pgxpool.Pool
}

type UserRepositoryInterface interface {
	NewUser() *User
	CreateUserTable() error
	GetUserByUsername(string) (*User, error)
	GetUserByEmail(string) (*User, error)
}

func (userRepo *UserRepository) InitialiseModels() { // only dealing with the userModel
	err := userRepo.CreateUserTable()
	if err != nil {
		slog.Error("Error Initialising User Repository")
		panic(err)
	}
}
