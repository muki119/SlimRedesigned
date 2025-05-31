package Models

import (
	"context"
	"github.com/jackc/pgx/v5"
	"v1/Config"
)

// type ModelErrors struct {
// 	groupTableError            error
// 	userTableError             error
// 	groupMessageTableError     error
// 	groupParticipantTableError error
// }

func InitialiseModels() {
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

	Config.DatabaseConnection.SendBatch(context.Background(), &batchTableCreation)
}
