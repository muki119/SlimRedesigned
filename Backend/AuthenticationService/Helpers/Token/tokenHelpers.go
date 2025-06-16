package Token

import (
	"crypto/rsa"
	"errors"
	"log"
)

var Token = createTokenService()

type HelperStruct struct {
	PrivateKey   *rsa.PrivateKey
	SymmetricKey []byte
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
)

func createTokenService() *HelperStruct {
	symmetricKey, err := getHMACSymmetricKey()
	if err != nil {
		log.Fatal("Error getting the symmetric key\n", err)
	}
	privateKey, err := getRSAPrivateKey()
	if err != nil {
		log.Fatal("error getting the private key", err)
	}
	return &HelperStruct{
		PrivateKey:   privateKey,
		SymmetricKey: symmetricKey,
	}
}

// make a refresh token  -- symmetric key  - only the auth server should verify ,
// make access token -- asymetric key - dont know what details maybe the userid
