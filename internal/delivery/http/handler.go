package http

import (
	"fmt"
	"net/http"
	"strconv"

	"AbdelrahmanDwedar/blogo/internal/usecase"
	"AbdelrahmanDwedar/blogo/pkg/response"
	"github.com/gorilla/mux"
)

// Handler aggregates all HTTP handlers
type Handler struct {
	UserHandler *UserHandler
	BlogHandler *BlogHandler
}

// NewHandler creates a new handler with all use cases
func NewHandler(userUC *usecase.UserUseCase, blogUC *usecase.BlogUseCase) *Handler {
	return &Handler{
		UserHandler: NewUserHandler(userUC),
		BlogHandler: NewBlogHandler(blogUC),
	}
}

// Ping is a health check endpoint
func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]string{"message": "pong"})
}

// Helper functions
func getIDFromPath(r *http.Request, key string) (int64, error) {
	vars := mux.Vars(r)
	idStr := vars[key]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid ID")
	}
	return id, nil
}

func getPaginationParams(r *http.Request) (limit, offset int) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit = 20 // default
	offset = 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	return limit, offset
}


