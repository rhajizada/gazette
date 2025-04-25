package oauth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func ExtractTokenFromHeaders(r *http.Request) (string, error) {
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		return "", errors.New("authorization header must be in format “Bearer {token}”")
	}
	return strings.TrimPrefix(auth, "Bearer "), nil
}

func VerifyToken(rawToken string, secret []byte) (*ApplicationClaims, error) {
	token, err := jwt.ParseWithClaims(rawToken, &ApplicationClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return secret, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(*ApplicationClaims)
	if !ok {
		return nil, errors.New("cannot parse claims")
	}

	return claims, nil
}
