package Config

import (
	"context"
	"github.com/redis/go-redis/v9"
)

var RedisContext = context.Background()
var RedisConnection = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})
