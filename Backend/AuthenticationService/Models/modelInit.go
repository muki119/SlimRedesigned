package Models

import (
	"github.com/jackc/pgx/v5"
	"log"
	"v1/Config"
)

func InitialiseModels() {
	log.Println("Initialising Models...")
	CreateGroupTable()
	CreateUserTable()
	CreateGroupMessageTable()
	CreateGroupParticipantTable()

	batchTableCreation := pgx.Batch{}
	batchTableCreation.Queue("CreateUserTable")
	batchTableCreation.Queue("CreateGroupTable")
	batchTableCreation.Queue("CreateGroupParticipantTable")
	batchTableCreation.Queue("CreateGroupMessageTable")
	batchTableCreation.Queue("CreateGroupMessageIndex")

	err := Config.DatabaseConnection.SendBatch(Config.DatabaseContext, &batchTableCreation).Close()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Initialised Models")
}
