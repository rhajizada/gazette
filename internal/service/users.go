package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// User is common model representing application user.
type User struct {
	ID            uuid.UUID `json:"id"`
	Sub           string    `json:"sub"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	CreatedAt     time.Time `json:"createdAt"`
	LastUpdatedAt time.Time `json:"lastUpdatedAt"`
}

// GetUserByID gets user by user ID.
func (s *Service) GetUserByID(ctx context.Context, userID uuid.UUID) (*User, error) {
	user, err := s.Repo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, NewError(
				fmt.Sprintf("user %s not found", userID),
				http.StatusNotFound,
			)
		} else {
			return nil, NewError(
				fmt.Sprintf("failed to fetch user %s", userID),
				http.StatusInternalServerError,
			)
		}
	}

	return &User{
		ID:            user.ID,
		Sub:           user.Sub,
		Name:          user.Name,
		Email:         user.Email,
		CreatedAt:     user.CreatedAt,
		LastUpdatedAt: user.CreatedAt,
	}, nil
}
