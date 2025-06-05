package Config

import (
	"context"
	"github.com/jackc/pgx/v5"
)

var DatabaseConnection *pgx.Conn
var DatabaseContext = context.Background()

func ConnectToDatabase() error {
	var databaseConfiguration, _ = pgx.ParseConfig("host=localhost port=5433 user=postgres dbname=slimDatabase")
	conn, err := pgx.Connect(DatabaseContext, databaseConfiguration.ConnString())
	DatabaseConnection = conn
	return err
}
