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

const (
	MinimumPasswordLength = 8
	MaximumPasswordLength = 96
)

var (
	ErrPasswordMinLength     = fmt.Errorf("password must be at least %d characters", MinimumPasswordLength)
	ErrInvalidHashedPassword = fmt.Errorf("invalid hashed password")
)

func HashPassword(plaintextPassword string) (string, error) {
	if len(plaintextPassword) < MinimumPasswordLength && len(plaintextPassword) > MaximumPasswordLength {
		return "", ErrPasswordMinLength
	}
	salt := make([]byte, 16) // the salt to use in bytes
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	passwordInBytes := []byte(plaintextPassword)
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

func ComparePassword(plaintextPassword string, hashedPassword string) (bool, error) {

	if len(plaintextPassword) < MinimumPasswordLength {
		return false, ErrPasswordMinLength
	}
	if len(hashedPassword) < MinimumPasswordLength {
		return false, ErrInvalidHashedPassword
	}
	splitPassword := strings.Split(hashedPassword, "$")
	if len(splitPassword) != 2 {
		return false, ErrInvalidHashedPassword
	}
	saltInBytes, err := base64.RawStdEncoding.DecodeString(splitPassword[0])
	if err != nil {
		return false, err
	}

	var hashedPasswordInBytes []byte
	hashedPasswordInBytes, err = base64.RawStdEncoding.DecodeString(splitPassword[1])
	if err != nil {
		return false, err
	}
	passwordInBytes := []byte(plaintextPassword)

	hashedIncomingPasswordInBytes, err := scrypt.Key(passwordInBytes, saltInBytes, n, r, p, keyLength)
	if err != nil {
		return false, err // return that its not equal and that
	}
	return bytes.Equal(hashedPasswordInBytes, hashedIncomingPasswordInBytes), nil

}
