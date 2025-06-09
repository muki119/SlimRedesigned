package TokenServices

import (
	"errors"
	"github.com/redis/go-redis/v9"
	"v1/Config"
)

func GetTokenService(tokenId string) string {
	Token, err := Config.RedisConnection.Get(Config.RedisContext, tokenId).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ""
		}
		return ""
	}
	return Token
}
