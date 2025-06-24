package Token

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"time"
	"v1/Config"
)

var (
	ErrTokenNotMapped = errors.New("token cannot be mapped")
)

func (*HelperStruct) BlockToken(token *jwt.Token) error {
	if token == nil {
		return ErrNoToken
	}
	return blockToken(token)
}
func (*HelperStruct) IsBlocklisted(token *jwt.Token) (bool, error) {
	if token == nil {
		return false, nil
	}
	tokenMap, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return true, ErrTokenNotMapped
	}
	tokenId, ok := tokenMap["jti"].(string)
	if !ok {
		return true, ErrNoTokenId
	}
	isBlocked, err := isBlocklisted(tokenId)
	if err != nil {
		return true, err
	}
	return isBlocked, nil
}

func isBlocklisted(tokenId string) (bool, error) {
	Token, err := Config.RedisConnection.Get(Config.RedisContext, tokenId).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, err
	}
	return Token != "", nil
}

func blockToken(RefreshToken *jwt.Token) error {
	tokenJti := RefreshToken.Claims.(jwt.MapClaims)["jti"]
	tokenExpiryDate := RefreshToken.Claims.(jwt.MapClaims)["exp"]
	ttl := time.Unix(int64(tokenExpiryDate.(float64)), 0).Sub(time.Now())
	status := Config.RedisConnection.Set(Config.RedisContext, tokenJti.(string), RefreshToken.Raw, ttl)
	if status.Err() != nil {
		return status.Err()
	}
	return nil
}
