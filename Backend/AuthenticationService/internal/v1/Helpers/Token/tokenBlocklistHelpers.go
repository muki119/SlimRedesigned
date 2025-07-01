package Token

import (
	"errors"
	"time"
	"v1/Config"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

var (
	ErrTokenNotMapped = errors.New("token cannot be mapped")
)

type Blocklist struct {
	Conn *redis.Client
}
type BlocklistInterface interface {
	BlockToken(token *jwt.Token) error
	IsBlocklisted(token *jwt.Token) (bool, error)
}

func (blocklist *Blocklist) BlockToken(token *jwt.Token) error {
	if token == nil {
		return ErrNoToken
	}
	return blocklist.blockToken(token)
}
func (blocklist *Blocklist) IsBlocklisted(token *jwt.Token) (bool, error) {
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
	isBlocked, err := blocklist.isBlocklisted(tokenId)
	if err != nil {
		return true, err
	}
	return isBlocked, nil
}

func (blocklist *Blocklist) isBlocklisted(tokenId string) (bool, error) {
	Token, err := blocklist.Conn.Get(Config.RedisContext, tokenId).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, err
	}
	return Token != "", nil
}

func (blocklist *Blocklist) blockToken(RefreshToken *jwt.Token) error {
	tokenJti := RefreshToken.Claims.(jwt.MapClaims)["jti"]
	tokenExpiryDate := RefreshToken.Claims.(jwt.MapClaims)["exp"]
	ttl := time.Until(time.Unix(int64(tokenExpiryDate.(float64)), 0))
	status := blocklist.Conn.Set(Config.RedisContext, tokenJti.(string), RefreshToken.Raw, ttl)
	if status.Err() != nil {
		return status.Err()
	}
	return nil
}
