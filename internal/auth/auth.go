package auth

import (
	"errors"
	"net/http"
	"strings"
)

// GetApiKey extracts an API key from the headers
// of an HTTP request
// Example:
// Authorization :ApiKey {insert apikey here}
func GetApiKey(headers http.Header) (string, error) {
    authorizationHeader := headers.Get("Authorization")
    if authorizationHeader == "" {
        return "", errors.New("no authentication info found")
    }

    authorization := strings.Split(authorizationHeader, " ")

    if len(authorization) != 2 {
        return "", errors.New("malformed auth header")
    }

    if authorization[0] != "ApiKey" {
        return "", errors.New("malformed first part of auth header")
    }

    return authorization[1], nil
}
