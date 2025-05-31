package Models

import (
	"context"
	"v1/Config"
)

type GroupMessages struct {
	Id        string `json:"id"`
	Group_Id  string `json:"group_id"`
	User_Id   string `json:"user_id"`
	Message   string `json:"message"`
	Date_Sent string `json:"date_sent"`
	Image     string `json:"image"`    // can be null or empty
	Reply_To  string `json:"reply_to"` // can be null
}

func CreateGroupMessageTable() {
	_, err := Config.DatabaseConnection.Prepare(context.Background(), "CreateGroupMessageTable", `
		CREATE TABLE IF NOT EXISTS group_messages (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			group_id UUID NOT NULL,
			user_id UUID NOT NULL,
			message TEXT NOT NULL,
			date_sent TIMESTAMP NOT NULL DEFAULT NOW(),
			image TEXT,
			reply_to UUID,
			FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (reply_to) REFERENCES group_messages(id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		panic(err)
	}
	// create index on group_id and timestamp
	_, err = Config.DatabaseConnection.Prepare(context.Background(), "CreateGroupMessageIndex", `
		CREATE INDEX IF NOT EXISTS group_messages_group_id_date_sent_index ON group_messages (group_id ASC,date_sent DESC)
	`)

	if err != nil {
		panic(err)
	}
}
