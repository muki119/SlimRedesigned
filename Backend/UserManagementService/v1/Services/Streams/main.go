package streams

import (
	"github.com/redis/go-redis/v9"
	"v1/Models"
)

type StreamServices struct {
	Connection     *redis.Client
	UserRepository *models.UserRepository
	// needs a redis streams connection
	// or a streams abstraction to only do the acrions we need
	// will also need a repository struct
	// only needs to delete and create profiles from redis stream messages currently
}
