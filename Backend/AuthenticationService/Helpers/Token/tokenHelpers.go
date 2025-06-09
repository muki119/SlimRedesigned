package Token

import (
	"crypto/rsa"
	"log"
)

var Token = createTokenService()

type HelperStruct struct {
	PrivateKey   *rsa.PrivateKey
	SymmetricKey []byte
}

func createTokenService() *HelperStruct {
	symmetricKey, err := getHMACSymmetricKey()
	if err != nil {
		log.Fatal(err)
	}
	privateKey, err := getRSAPrivateKey()
	if err != nil {
		log.Fatal(err)
	}
	return &HelperStruct{
		PrivateKey:   privateKey,
		SymmetricKey: symmetricKey,
	}
}

// make a refresh token  -- symmetric key  - only the auth server should verify ,
// make access token -- asymetric key - dont know what details maybe the userid
