package models

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
}

type UserRepository struct {
	DatabaseConn *pgxpool.Pool
}
