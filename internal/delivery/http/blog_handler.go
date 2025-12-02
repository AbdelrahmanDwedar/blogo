package http

import (
	"encoding/json"
	"net/http"

	"AbdelrahmanDwedar/blogo/internal/domain/entity"
	"AbdelrahmanDwedar/blogo/internal/usecase"
	"AbdelrahmanDwedar/blogo/pkg/auth"
	"AbdelrahmanDwedar/blogo/pkg/response"
)

// BlogHandler handles blog-related HTTP requests
type BlogHandler struct {
	blogUC *usecase.BlogUseCase
}

// NewBlogHandler creates a new blog handler
func NewBlogHandler(blogUC *usecase.BlogUseCase) *BlogHandler {
	return &BlogHandler{blogUC: blogUC}
}

// CreateBlog creates a new blog post
func (h *BlogHandler) CreateBlog(w http.ResponseWriter, r *http.Request) {
	claims, err := auth.GetUserFromContext(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Body        string `json:"body"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	blog, err := h.blogUC.CreateBlog(req.Title, req.Description, req.Body, claims.UserID)
	if err != nil {
		if err == entity.ErrInvalidTitle || err == entity.ErrInvalidBody {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "Failed to create blog")
		return
	}

	response.Created(w, blog)
}

// GetBlogs retrieves all blogs
func (h *BlogHandler) GetBlogs(w http.ResponseWriter, r *http.Request) {
	limit, offset := getPaginationParams(r)

	blogs, err := h.blogUC.GetAllBlogs(limit, offset)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to get blogs")
		return
	}

	response.Success(w, map[string]interface{}{
		"blogs":  blogs,
		"limit":  limit,
		"offset": offset,
	})
}

// GetBlog retrieves a single blog by ID
func (h *BlogHandler) GetBlog(w http.ResponseWriter, r *http.Request) {
	blogID, err := getIDFromPath(r, "id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid blog ID")
		return
	}

	blog, err := h.blogUC.GetBlogByID(blogID)
	if err != nil {
		if err == entity.ErrBlogNotFound {
			response.Error(w, http.StatusNotFound, "Blog not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Failed to get blog")
		return
	}

	response.Success(w, blog)
}

// UpdateBlog updates a blog post
func (h *BlogHandler) UpdateBlog(w http.ResponseWriter, r *http.Request) {
	claims, err := auth.GetUserFromContext(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	blogID, err := getIDFromPath(r, "id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid blog ID")
		return
	}

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Body        string `json:"body"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	blog, err := h.blogUC.UpdateBlog(blogID, req.Title, req.Description, req.Body, claims.UserID)
	if err != nil {
		if err == entity.ErrNotBlogOwner {
			response.Error(w, http.StatusForbidden, "You can only edit your own blogs")
			return
		}
		if err == entity.ErrBlogNotFound {
			response.Error(w, http.StatusNotFound, "Blog not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Failed to update blog")
		return
	}

	response.Success(w, blog)
}

// DeleteBlog deletes a blog post
func (h *BlogHandler) DeleteBlog(w http.ResponseWriter, r *http.Request) {
	claims, err := auth.GetUserFromContext(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	blogID, err := getIDFromPath(r, "id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid blog ID")
		return
	}

	err = h.blogUC.DeleteBlog(blogID, claims.UserID)
	if err != nil {
		if err == entity.ErrNotBlogOwner {
			response.Error(w, http.StatusForbidden, "You can only delete your own blogs")
			return
		}
		if err == entity.ErrBlogNotFound {
			response.Error(w, http.StatusNotFound, "Blog not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Failed to delete blog")
		return
	}

	response.Success(w, map[string]string{
		"message": "Blog deleted successfully",
	})
}

// LikeBlog likes or unlikes a blog
func (h *BlogHandler) LikeBlog(w http.ResponseWriter, r *http.Request) {
	claims, err := auth.GetUserFromContext(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	blogID, err := getIDFromPath(r, "id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid blog ID")
		return
	}

	var req struct {
		Action string `json:"action"` // "like" or "unlike"
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.Action == "like" {
		err = h.blogUC.LikeBlog(blogID, claims.UserID)
	} else if req.Action == "unlike" {
		err = h.blogUC.UnlikeBlog(blogID, claims.UserID)
	} else {
		response.Error(w, http.StatusBadRequest, "Invalid action. Use 'like' or 'unlike'")
		return
	}

	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to "+req.Action+" blog")
		return
	}

	response.Success(w, map[string]string{
		"message": "Successfully " + req.Action + "d blog",
	})
}

// GetBlogLikes retrieves users who liked a blog
func (h *BlogHandler) GetBlogLikes(w http.ResponseWriter, r *http.Request) {
	blogID, err := getIDFromPath(r, "id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid blog ID")
		return
	}

	limit, offset := getPaginationParams(r)

	likes, err := h.blogUC.GetBlogLikes(blogID, limit, offset)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to get likes")
		return
	}

	response.Success(w, map[string]interface{}{
		"likes":  likes,
		"limit":  limit,
		"offset": offset,
	})
}


