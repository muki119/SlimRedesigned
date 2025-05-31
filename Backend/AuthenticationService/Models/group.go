package Models

import (
	"context"
	"fmt"
	"v1/Config"
)

type Group struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Date_Created string `json:"created_at"`
	Created_By   string `json:"created_by"`
	Admin        string `json:"admin"`
}

func CreateGroupTable() error {
	out, err := Config.DatabaseConnection.Prepare(context.Background(), "CreateGroupTable", `
		CREATE TABLE IF NOT EXISTS groups (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(100) NOT NULL UNIQUE,
			description VARCHAR(255) NOT NULL,
			date_created TIMESTAMP NOT NULL DEFAULT NOW(),
			created_by UUID NOT NULL,
			admin UUID NOT NULL,
			FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE NO ACTION,
			FOREIGN KEY (admin) REFERENCES users(id) ON DELETE NO ACTION
		)
	`)

	fmt.Println(out)
	return err
}
