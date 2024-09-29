package auth

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 1)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}
