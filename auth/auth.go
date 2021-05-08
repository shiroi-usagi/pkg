package auth

import (
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
)

// ErrNoBearerToken is returned by BearerToken when no authorization header
// is set or the authorization header is not a bearer token.
var ErrNoBearerToken = errors.New("token: no bearer token in request header")

// BearerToken returns the bearer token from the request when available.
// It returns an ErrNoBearerToken error when missing or a hex error when the
// token is not a valid hex.
func BearerToken(r *http.Request) ([]byte, error) {
	authorization := r.Header.Get("Authorization")
	if !strings.HasPrefix(authorization, "Bearer ") {
		return nil, ErrNoBearerToken
	}
	return hex.DecodeString(strings.TrimPrefix(authorization, "Bearer "))
}
