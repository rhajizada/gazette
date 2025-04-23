package oauth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Name   string   `json:"name"`
	Email  string   `json:"email"`
	Sub    string   `json:"sub"`
	Groups []string `json:"groups"`
}

func (c *Claims) GetJwtMap(expiration time.Duration) jwt.MapClaims {
	return jwt.MapClaims{
		"name":   c.Name,
		"email":  c.Email,
		"sub":    c.Sub,
		"groups": c.Groups,
		"exp":    time.Now().Add(expiration).Unix(),
	}
}
