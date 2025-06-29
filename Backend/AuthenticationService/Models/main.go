package Models

import (
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	Db *pgxpool.Pool
}
type UserRepositoryInterface interface {
	NewUser() *User
	CreateUserTable()
	GetUserByUsername(string) (*User, error)
	GetUserByEmail(string) (*User, error)
}

func (userRepo *UserRepository) InitialiseModels() { // only dealing with the userModel
	err := userRepo.CreateUserTable()
	if err != nil {
		fmt.Println("Error Initialising User Repository")
		panic(err)
	}
}
