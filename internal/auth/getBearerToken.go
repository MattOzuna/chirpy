package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	token := headers.Get("Authorization")
	if token == "" {
		err := fmt.Errorf("no 'Authorization' header")
		return "", err
	}
	token = strings.TrimSpace(token)
	return token, nil
}
