package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestValidateJWT(t *testing.T) {

	testId := uuid.New()

	tests := []struct {
		name      string
		expiresIn time.Duration
		userId    uuid.UUID
		secret    string
		expected  string
	}{
		{
			name:      "normal test",
			expiresIn: 24 * time.Hour,
			userId:    testId,
			secret:    "super secret",
		},
		// {
		// 	name:      "expired token",
		// 	expiresIn: 1 * time.Second,
		// 	userId:    testId,
		// 	secret:    "super secret",
		// },
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			token, _ := MakeJWT(tc.userId, tc.secret, tc.expiresIn)
			// time.Sleep(2 * time.Second)
			user, err := ValidateJWT(token, "super secret")
			t.Logf("\nuserID: %v", user)
			if err != nil {
				t.Fatalf("Validate test failed: %v", err)
			}
		})
	}
}
