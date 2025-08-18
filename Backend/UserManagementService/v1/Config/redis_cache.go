package config

import "github.com/redis/go-redis/v9"

// will cache user profiles with a time to live of around 2-3 days
// when accessed again a profiles time to live will be refreshed
// on update of a user from authentication service , will update users info
type RedisCache struct {
	Addr     string
	Password string
	DB       int
}

func (rc *RedisCache) Connect() (*redis.Client, error) {

	return nil, nil
}
