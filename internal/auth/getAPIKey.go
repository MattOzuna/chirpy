package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	key := headers.Get("Authorization")
	if key == "" {
		err := fmt.Errorf("no 'Authorization' header")
		return "", err
	}
	key = strings.TrimPrefix(key, "ApiKey ")
	return key, nil
}
