package helpers

import (
	"errors"
	"os"
)

func GetEnv(key string) (string, error) {
	if value, exists := os.LookupEnv(key); exists {
		return value, nil
	}
	return "", errors.New(key + " is required")
}
