package oauth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type ProviderClaims struct {
	Name   string   `json:"name"`
	Email  string   `json:"email"`
	Sub    string   `json:"sub"`
	Groups []string `json:"groups"`
}

type ApplicationClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Name   string    `json:"name"`
	Email  string    `json:"email"`
	Sub    string    `json:"sub"`
	Groups []string  `json:"groups"`
	jwt.RegisteredClaims
}

func (c *ProviderClaims) GetAppClaims(userID uuid.UUID, expiration time.Duration) jwt.MapClaims {
	return jwt.MapClaims{
		"user_id": userID.String(),
		"name":    c.Name,
		"email":   c.Email,
		"sub":     c.Sub,
		"groups":  c.Groups,
		"iss":     "gazette",
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(expiration).Unix(),
	}
}
