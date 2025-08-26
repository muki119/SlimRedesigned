package streams

import (
	"v1/Models"
)

type StreamServices struct {
	UserRepository *models.UserRepository
	// needs a redis streams connection
	// or a streams abstraction to only do the acrions we need
	// will also need a repository struct
	// only needs to delete and create profiles from redis stream messages currently
}
