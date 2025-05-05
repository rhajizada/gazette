package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/rhajizada/gazette/internal/service"

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
		var serviceErr service.ServiceError
		if errors.As(err, &serviceErr) {
			http.Error(w, serviceErr.Error(), int(serviceErr.Code))
			return
		} else {
			http.Error(w, "failed to fetch user", http.StatusBadRequest)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
