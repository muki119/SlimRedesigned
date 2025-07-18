package Token

import (
	"crypto/rsa"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"log"
)

type Helper struct {
	PrivateKey   *rsa.PrivateKey
	SymmetricKey []byte
	Blocklist    *Blocklist
}

type helperTokenInterface interface {
	CreateAccessToken(string, string) (string, error)
	CreateLoginRefreshToken(string) (string, error)
	ParseRefreshToken(string) (*jwt.Token, error)
	CreateRefreshTokenFromClaims(jwt.Claims) (string, error)
}

var (
	ErrNoToken        = errors.New("no token")
	ErrNoIssuer       = errors.New("no issuer")
	ErrInvalidToken   = errors.New("invalid token")
	ErrNoUserId       = errors.New("no user id")
	ErrNoClaims       = errors.New("no claims")
	ErrNoExpiry       = errors.New("no expiry")
	ErrNoPrivateKey   = errors.New("no private key")
	ErrNoSymmetricKey = errors.New("no symmetric key")
	ErrNoTokenId      = errors.New("no token id")
)

type HelperTokenConfig struct {
	//the environment key for the private key
	PrivateKey string
	//the environment key for the private key
	SecretKey string
}

func (config *HelperTokenConfig) CreateTokenService(db *redis.Client) *Helper {
	secretKey, err := GetHMACSymmetricKey(config.SecretKey)
	if err != nil {
		log.Fatal("Error getting the symmetric key\n", err)
	}
	privateKey, err := GetRSAPrivateKey(config.PrivateKey)
	if err != nil {
		log.Fatal("error getting the private key", err)
	}
	return &Helper{
		PrivateKey:   privateKey,
		SymmetricKey: secretKey,
		Blocklist: &Blocklist{
			Conn: db,
		},
	}
}

// make a refresh token  -- symmetric key  - only the auth server should verify ,
// make access token -- asymetric key - dont know what details maybe the userid
