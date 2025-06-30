package Config

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	DatabaseConnection *pgxpool.Pool
	DatabaseContext    = context.Background()
)

type PGDatabase struct {
	Host string
	Port string
	User string
	Name string
}

func (db *PGDatabase) ConnectToDatabase() (*pgxpool.Pool, error) {

	connectionString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s", db.Host, db.Port, db.User, db.Name)
	var databaseConfiguration, err = pgx.ParseConfig(connectionString)
	if err != nil {
		return nil, err
	}
	conn, err := pgxpool.New(DatabaseContext, databaseConfiguration.ConnString())
	if err != nil {
		return nil, err
	}

	return conn, nil
}
