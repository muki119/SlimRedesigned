package Config

import (
	"context"
	"github.com/jackc/pgx/v5"
)

var DatabaseConnection *pgx.Conn

func ConnectToDatabase() error {
	var databaseConfiguration, _ = pgx.ParseConfig("host=localhost port=5433 user=postgres dbname=slimDatabase")
	conn, err := pgx.Connect(context.Background(), databaseConfiguration.ConnString())
	DatabaseConnection = conn
	return err
}
