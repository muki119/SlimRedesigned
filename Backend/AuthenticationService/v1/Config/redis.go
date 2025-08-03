package Config

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisBlocklistConfig struct {
	Addr     string
	Password string
	DB       int
}
type RedisClient interface {
	Connect() error
}

var RedisContext = context.Background()

func (config *RedisBlocklistConfig) ConnectToDatabase() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})
}
