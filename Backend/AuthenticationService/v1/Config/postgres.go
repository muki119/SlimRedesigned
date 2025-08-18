package Config

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	DatabaseContext = context.Background()
)

type PGDatabase struct {
	Host     string
	Port     string
	User     string
	Name     string
	Password string
	Timeout  string
}
type PGDatabaseInterface interface {
	ConnectToDatabase() (*pgxpool.Pool, error)
}

func (db *PGDatabase) ConnectToDatabase() (*pgxpool.Pool, error) {
	if db.Host == "" || db.Port == "" || db.User == "" || db.Name == "" {
		return nil, fmt.Errorf("database connection parameters are not set")
	}
	connectionString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s connect_timeout=%s", db.Host, db.Port, db.User, db.Name, db.Password, db.Timeout)
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
