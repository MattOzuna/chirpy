package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	expiresIn := 24 * time.Hour
	userID := uuid.New()
	secret := "superSecret"

	_, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("makeJWT test failed: %v", err)
	}
}
