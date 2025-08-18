package Utils

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"
)

type ErrNoEnvFound struct {
	key string
}

func (e ErrNoEnvFound) Error() string {
	return fmt.Sprintf("no env found for key %s", e.key)
}

var ErrCannotConvert = func(key string) error {
	return fmt.Errorf("cannot convert key %s", key)
}

func MustGetEnv(key string) string {
	variable := os.Getenv(key)
	if variable == "" {
		slogFatal(ErrNoEnvFound{key})
	}
	return variable
}

func MustGetEnvInt(key string) int {
	variable := os.Getenv(key)
	if variable == "" {
		slogFatal(ErrNoEnvFound{key})
	}
	value, err := strconv.Atoi(variable)
	if err != nil {
		slogFatal(ErrCannotConvert(key))
	}
	return value

}

func slogFatal(err error) {
	if err != nil {
		slog.Log(context.Background(), slog.LevelError, "Fatal error", "error", err)
		os.Exit(1)
	}
}
