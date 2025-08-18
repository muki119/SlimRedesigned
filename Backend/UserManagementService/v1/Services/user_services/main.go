package userservices

import (
	models "v1/Models"
)

type UserServices struct {
	UserRepository *models.UserRepository
	// most likely needs a repository struct
}
