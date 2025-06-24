package Models

import (
	"log"
	"v1/Config"

	"github.com/jackc/pgx/v5"
)

func InitialiseModels() { // only dealing with the userModel
	log.Println("Initialising Models...")
	CreateUserTable()
	batchTableCreation := pgx.Batch{}
	batchTableCreation.Queue("CreateUserTable")
	err := Config.DatabaseConnection.SendBatch(Config.DatabaseContext, &batchTableCreation).Close()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Initialised Models")
}
