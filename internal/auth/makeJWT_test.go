package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	expiresIn := 24 * time.Hour
	userID := uuid.New()
	secret := "super secret"

	token, err := MakeJWT(userID, secret, expiresIn)
	t.Logf("\n token: %v", token)
	if err != nil {
		t.Fatalf("makeJWT test failed: %v", err)
	}
}
