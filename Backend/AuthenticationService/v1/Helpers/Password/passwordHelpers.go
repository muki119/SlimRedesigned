package Password

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/scrypt"
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
	splitHashedPassword := strings.Split(hashedPassword, "$") // split the hashed password to get hash and password ciphertext
	if len(splitHashedPassword) != 2 {
		return false, ErrInvalidHashedPassword
	}
	bytesSalt, err := base64.RawStdEncoding.DecodeString(splitHashedPassword[0])
	if err != nil {
		return false, err // return that it's not equal and that theres an error
	}

	bytesHashedPasswordValue, err := base64.RawStdEncoding.DecodeString(splitHashedPassword[1]) // the stored hashed passwo
	if err != nil {
		return false, err
	}

	bytesPlaintextPassword := []byte(plaintextPassword)
	bytesHashedPlaintextPassword, err := scrypt.Key(bytesPlaintextPassword, bytesSalt, n, r, p, keyLength) // hashing the plaintext password
	if err != nil {
		return false, err // return that it's not equal and that theres an error
	}
	return subtle.ConstantTimeCompare(bytesHashedPlaintextPassword, bytesHashedPasswordValue) == 1, nil

}
