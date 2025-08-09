package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/rhajizada/gazette/internal/middleware"
	"github.com/rhajizada/gazette/internal/repository"
	"github.com/rhajizada/gazette/internal/service"
)

// returns paginated list of recent items from subscribed user feeds.
// @Summary      List subscribed items
// @Description  Retrieves subscribed items for user, paginated.
// @Tags         Subscribed
// @Param        limit   query     int32  true   "Max number of items"
// @Param        offset  query     int32  true   "Number of items to skip"
// @Success      200     {object}  service.ListItemsResponse
// @Failure      400     {object}  string
// @Failure      500     {object}  string
// @Security     BearerAuth
// @Router       /api/subscribed [get]
func (h *Handler) ListUserSubscribedItems(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserClaims(r).UserID

	params, err := getPageParams(r.URL.Query())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var resp *service.ListItemsResponse
	resp, err = h.Service.ListUserSubscribedItems(r.Context(),
		repository.ListSubscribedItemsByUserParams{
			UserID: userID,
			Limit:  params.Limit,
			Offset: params.Offset,
		})
	if err != nil {
		var serviceErr service.ServiceError
		if errors.As(err, &serviceErr) {
			http.Error(w, serviceErr.Error(), int(serviceErr.Code))
			return
		} else {
			http.Error(w, "failed to list subscribed items", http.StatusBadRequest)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
