package Utils

import (
	"errors"
	"os"
)

func Getenv(key string) string {
	variable := os.Getenv(key)
	var ErrNoEnvVariable = errors.New("No such environment variable " + key)
	if variable == "" {
		panic(ErrNoEnvVariable)
	}
	return variable
}
