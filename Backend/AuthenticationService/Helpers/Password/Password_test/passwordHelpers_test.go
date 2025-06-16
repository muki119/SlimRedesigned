package Password_test

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"testing"
	"v1/Helpers/Password"
)

var randomArrayOfStrings = generateRandomArrayOfStrings(20)

func TestHashPassword(t *testing.T) {
	t.Run("ReturnsHashWhenGivenInput", func(t *testing.T) {
		passwordToBeHashed := randomArrayOfStrings[0]
		hashedPassword, err := Password.HashPassword(passwordToBeHashed)
		if err != nil {
			t.Error(err)
		}
		if hashedPassword == "" {
			t.Error("HashedPassword should not be empty")
		}
	})
	t.Run("ReturnsErrorWhenGivenNoInput", func(t *testing.T) {
		_, err := Password.HashPassword("")
		if err == nil {
			t.Error("HashedPassword should return an error")
		}
	})
	t.Run("ReturnsErrorWhenGivenInvalidLength", func(t *testing.T) { // should be more than 6
		_, err := Password.HashPassword("invalid")
		if err == nil {
			t.Error("HashedPassword should return an error")
		}
		if !errors.Is(err, Password.ErrPasswordMinLength) {
			t.Error("HashedPassword should return an PasswordMinLength Error")
		}
	})
}
func TestComparePassword(t *testing.T) {

	plainTextPassword := randomArrayOfStrings[0]
	hashedPassword, err := Password.HashPassword(plainTextPassword)
	if err != nil {
		t.Error(err)
	}
	t.Run("ReturnsTrueWhenPasswordsAreTheSame", func(t *testing.T) {
		comparedPassword, err := Password.ComparePassword(plainTextPassword, hashedPassword)
		if err != nil {
			t.Error(err)
		}
		if !comparedPassword {
			t.Error("Password should be the same")
		}
	})
	t.Run("ReturnsFalseWhenPasswordsAreNotTheSame", func(t *testing.T) {
		comparedPassword, err := Password.ComparePassword(hashedPassword, plainTextPassword)
		if err != nil {
			t.Error(err)
		}
		if comparedPassword {
			t.Error("Password should not be the same")
		}
	})

	t.Run("ReturnsErrorWhenNoPlaintextPassword", func(t *testing.T) {
		_, err := Password.ComparePassword("", hashedPassword)
		if err == nil {
			t.Error("HashedPassword should return an error")
		}
		if !errors.Is(err, Password.ErrPasswordMinLength) {
			t.Error("Password should return an PasswordMinLength error")
		}
	})
	t.Run("ReturnsErrorWhenNoHashedPassword", func(t *testing.T) {
		_, err := Password.ComparePassword(plainTextPassword, "")
		if err == nil {
			t.Error("HashedPassword should return an error")
		}
		if !errors.Is(err, Password.ErrInvalidHashedPassword) {
			t.Error("HashedPassword should return InvalidHashedPassword error")
		}
	})

}

func generateRandomArrayOfStrings(length int) []string {
	arrOfStrings := make([]string, length)
	for i := 0; i < length; i++ {
		randomEightByteString := make([]byte, 8)
		_, err := rand.Read(randomEightByteString)
		if err != nil {
			return nil
		}
		arrOfStrings[i] = hex.EncodeToString(randomEightByteString)
		fmt.Println(arrOfStrings[i])
	}
	return arrOfStrings
}
