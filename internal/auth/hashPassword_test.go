package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "password"
	hash, err := HashPassword(password)
	if hash == "password" || err != nil {
		t.Fatalf("hash test failed: %v", err)
	}
}
