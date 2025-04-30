package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	_ "github.com/rhajizada/gazette/internal/service"

	"github.com/rhajizada/gazette/internal/middleware"
)

// GetUser returns currently logged in user.
// @Summary      Get user
// @Description  Retrieves currently logged in user.
// @Tags         Users
// @Success      200     {object}  service.User
// @Failure      400     {object}  string
// @Failure      404     {object}  string
// @Failure      500     {object}  string
// @Security     BearerAuth
// @Router       /api/user [get]
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserClaims(r).UserID
	user, err := h.Service.GetUserByID(r.Context(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "user not found", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("failed user item: %v", err), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
