package auth

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestCheckPassword(t *testing.T) {
	password := "password"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 1)
	if err != nil {
		t.Fatal(err)
	}

	if err := CheckPassword(password, string(hashedPassword)); err != nil {
		t.Fatalf("checkPassword test failed: %v", err)
	}
}
