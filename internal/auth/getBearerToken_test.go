package auth

import (
	"net/http"
	"testing"
)

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name     string
		header   http.Header
		expected string
	}{
		{
			name: "valid token",
			header: http.Header{
				"Authorization": {"test_string"},
			},
			expected: "test_string",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			token, err := GetBearerToken(tc.header)
			t.Log(tc.name)
			if err != nil {
				t.Fatal(err)
			}
			if token != tc.expected {
				t.Fatal("token doesn't match expected value")
			}
		})
	}
}
