package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/rhajizada/gazette/internal/middleware"
	"github.com/rhajizada/gazette/internal/repository"
	"github.com/rhajizada/gazette/internal/service"
)

// ListCategories returns distinct categories.
// @Summary      List categories
// @Description  Retrieves paginated list of distinct categories.
// @Tags         Categories
// @Param        limit   query     int32  true   "Max number of categories"
// @Param        offset  query     int32  true   "Number of categories to skip"
// @Success      200     {object}  service.ListCategoriesResponse
// @Failure      400     {object}  string
// @Failure      500     {object}  string
// @Security     BearerAuth
// @Router       /api/categories [get]
func (h *Handler) ListCategories(w http.ResponseWriter, r *http.Request) {
	_ = middleware.GetUserClaims(r).UserID

	params, err := getPageParams(r.URL.Query())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var resp *service.ListCategoriesResponse

	resp, err = h.Service.ListCategories(r.Context(),
		repository.ListDistinctCategoriesParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		})
	if err != nil {
		var serviceErr service.ServiceError
		if errors.As(err, &serviceErr) {
			http.Error(w, serviceErr.Error(), int(serviceErr.Code))
			return
		} else {
			http.Error(w, "failed to list categories", http.StatusBadRequest)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// ListItemsByCategories returns paginated list of items in given categories.
// @Summary      List items in the categories.
// @Description  Retrieves items in the categories, including like status.
// @Tags         Categories
// @Param        names   query     []string   true   "Category names"  collectionFormat(multi)
// @Param        limit        query     int32      true   "Max number of items"
// @Param        offset       query     int32      true   "Number of items to skip"
// @Success      200          {object}  service.ListItemsResponse
// @Failure      400          {object}  string
// @Failure      500          {object}  string
// @Security     BearerAuth
// @Router       /api/categories/items [get]
func (h *Handler) ListItemsByCategories(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserClaims(r).UserID

	// collect all "categories" params into a []string
	categories := r.URL.Query()["names"]
	if len(categories) == 0 {
		http.Error(w, "bad input: missing categories", http.StatusBadRequest)
		return
	}

	params, err := getPageParams(r.URL.Query())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.Service.ListCategoryItems(r.Context(),
		repository.ListItemsByCategoryForUserParams{
			Categories: categories,
			UserID:     userID,
			Limit:      params.Limit,
			Offset:     params.Offset,
		},
	)
	if err != nil {
		var serviceErr service.ServiceError
		if errors.As(err, &serviceErr) {
			http.Error(w, serviceErr.Error(), int(serviceErr.Code))
		} else {
			http.Error(w,
				fmt.Sprintf("failed to list items in categories %v", categories),
				http.StatusBadRequest,
			)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
