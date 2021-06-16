package interserviceclient

import (
	"fmt"
	"net/http"
	"strings"
)

// ExtractBearerToken gets a bearer token from an Authorization header.
// This is expected to contain a Firebase idToken prefixed with "Bearer "
func ExtractBearerToken(r *http.Request) (string, error) {
	return ExtractToken(r, "Authorization", "Bearer")
}

// ExtractToken extracts a token with the specified prefix from the specified header
func ExtractToken(r *http.Request, header string, prefix string) (string, error) {
	if r == nil {
		return "", fmt.Errorf("nil request")
	}
	if r.Header == nil {
		return "", fmt.Errorf("no headers, can't extract bearer token")
	}
	authHeader := r.Header.Get(header)
	if authHeader == "" {
		return "", fmt.Errorf("expected an `%s` request header", header)
	}
	if !strings.HasPrefix(authHeader, prefix) {
		return "", fmt.Errorf("the `Authorization` header contents should start with `Bearer`")
	}
	tokenOnly := strings.TrimSpace(strings.TrimPrefix(authHeader, prefix))
	return tokenOnly, nil
}
