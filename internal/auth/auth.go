package auth

import (
	"errors"
	"net/http"
	"strings"
)

type AuthConfig struct {
	ApiKey string
}

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	key, ok := strings.CutPrefix(authHeader, "ApiKey ")
	if !ok {
		return "", errors.New("malformed Authorization header")
	}
	return key, nil
}

func (cfg *AuthConfig) AuthenticateAPIKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := GetAPIKey(r.Header)
		if err != nil || apiKey != cfg.ApiKey {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Not authorized"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
