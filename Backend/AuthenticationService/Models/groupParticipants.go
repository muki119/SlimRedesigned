package Models

import (
	"context"
	"v1/Config"
)

type GroupParticipants struct {
	User_Id     string `json:"user_id"`
	Group_Id    string `json:"group_id"`
	Date_Joined string `json:"date_joined"`
}

func CreateGroupParticipantTable() { // Initial function to create the group participants table
	_, err := Config.DatabaseConnection.Prepare(context.Background(), "CreateGroupParticipantTable", `
		CREATE TABLE IF NOT EXISTS group_participants (
			user_id UUID NOT NULL,
			group_id UUID NOT NULL,
			date_joined TIMESTAMP NOT NULL DEFAULT NOW(),
			PRIMARY KEY (group_id,user_id),
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		panic(err)
	}
}
