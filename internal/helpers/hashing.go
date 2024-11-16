package helpers

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

var (
	PASSWORD_SALT = os.Getenv("PASSWORD_SALT")
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(fmt.Sprintf("%s;%s", PASSWORD_SALT, password)), 12)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(fmt.Sprintf("%s;%s", PASSWORD_SALT, password)))
	return err == nil
}
