package Password

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/scrypt"
	"strings"
)

const (
	n         = 16384
	r         = 8
	p         = 1
	keyLength = 32
)

func HashPassword(password string) (string, error) {
	salt := make([]byte, 16) // the salt to use in bytes
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	passwordInBytes := []byte(password)
	var hashedPassword []byte
	hashedPassword, err = scrypt.Key(passwordInBytes, salt, n, r, p, keyLength)
	if err != nil {
		return "", err
	}
	base64encodedPassword := base64.RawStdEncoding.EncodeToString(hashedPassword)
	base64encodedSalt := base64.RawStdEncoding.EncodeToString(salt)
	saltAndPasswordString := fmt.Sprintf("%s$%s", base64encodedSalt, base64encodedPassword) // store as concat string with $ delimiter
	return saltAndPasswordString, nil
}

func ComparePassword(password string, hashedPassword string) (bool, error) {
	splitPassword := strings.Split(hashedPassword, "$")
	saltInBytes, _ := base64.RawStdEncoding.DecodeString(splitPassword[0])

	var hashedPasswordInBytes []byte
	hashedPasswordInBytes, _ = base64.RawStdEncoding.DecodeString(splitPassword[1])
	passwordInBytes := []byte(password)

	hashedIncomingPasswordInBytes, err := scrypt.Key(passwordInBytes, saltInBytes, n, r, p, keyLength)
	if err != nil {
		return false, err // return that its not equal and that
	}
	return bytes.Equal(hashedPasswordInBytes, hashedIncomingPasswordInBytes), nil

}
