package Model_test

import (
	"AuthenticationService/v1/Config"
	"AuthenticationService/v1/Models"
	"os"
	"testing"
)

// test if save user is all good
func TestSaveUser(t *testing.T) {
	testDBConfig := Config.PGDatabase{
		Host:    os.Getenv("DB_HOST"),
		Port:    os.Getenv("DB_PORT"),
		User:    os.Getenv("DB_USER"),
		Name:    os.Getenv("DB_TEST_NAME"),
		Timeout: os.Getenv("DB_TIMEOUT"),
	}
	testDb, err := testDBConfig.ConnectToDatabase()
	testUserRepo := Models.UserRepository{Db: testDb}
	if err != nil {
		t.Error(err)
	}
	testUserRepo.InitialiseModels()
	defer testUserRepo.Db.Close()

	t.Run("SaveUserShouldSave", func(t *testing.T) {
		// create a regular test account and try to make an account, should work, on cleanup should delete a test account
		validUser := testUserRepo.NewUser()

		validUser.Forename = "Stannis"
		validUser.Surname = "Baratheon"
		validUser.Username = "stannisB123"
		validUser.Email = "stannisB123@gotmail.com"
		validUser.Password = "ih8ramseyBolton"
		validUser.DateOfBirth = "2025-06-26T14:41:43.227Z"
		validUser.Role = "USER"

		err := validUser.SaveUser()
		if err != nil {
			t.Error("Should be able to be able to save user", err)
		}
		t.Cleanup(func() {
			deleteUser, err := testUserRepo.GetUserByUsername("stannisB123")
			if err != nil {
				t.Error("Should be able to be able to retrieve user", err)
			}
			err = deleteUser.Delete()
			if err != nil {
				return
			}
		})
	})
	t.Run("SaveUserShouldReturnErr", func(t *testing.T) {
		// test should create some bad info and try to save it ,
		invalidUser := testUserRepo.NewUser()

		invalidUser.Forename = "Stannis"
		invalidUser.Surname = "thirtyonecharachtersmightbeenoughtocauseacodeerrorbutthisismore"
		invalidUser.Username = "stannisB123"
		invalidUser.Email = "stannisB123@gotmail.com"
		invalidUser.Password = "ih8ramseyBolton"
		invalidUser.DateOfBirth = "2025-06-26T14:41:43.227Z"
		invalidUser.Role = "USER"

		err := invalidUser.SaveUser()
		if err == nil {
			t.Error("shouldn't be able to save user")
			user, err := testUserRepo.GetUserByUsername(invalidUser.Username)
			if err != nil {
				t.Error(err)
				return
			}
			err = user.Delete()
			if err != nil {
				t.Error(err)
				return
			}
		}

	})
}
