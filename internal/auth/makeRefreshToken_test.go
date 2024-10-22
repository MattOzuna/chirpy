package auth

import (
	"testing"
)

func TestMakeRefreshTokenTest(t *testing.T) {
	_, err := MakeRefreshToken()
	if err != nil {
		t.Fatalf("makeRefreshToken test failed: %v", err)
	}
}
