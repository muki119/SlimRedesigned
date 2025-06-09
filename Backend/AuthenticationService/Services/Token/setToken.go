package TokenServices

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
	"v1/Config"
	//"v1/Config"
)

func SetTokenService(RefreshToken *jwt.Token) error {
	tokenJti := RefreshToken.Claims.(jwt.MapClaims)["jti"]
	tokenExpiryDate := RefreshToken.Claims.(jwt.MapClaims)["exp"]
	ttl := time.Unix(int64(tokenExpiryDate.(float64)), 0).Sub(time.Now())
	status := Config.RedisConnection.Set(Config.RedisContext, tokenJti.(string), RefreshToken.Raw, ttl)
	if status.Err() != nil {
		return status.Err()
	}
	return nil
}
