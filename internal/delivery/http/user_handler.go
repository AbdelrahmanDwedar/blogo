package http

import (
	"encoding/json"
	"net/http"

	"AbdelrahmanDwedar/blogo/internal/domain/entity"
	"AbdelrahmanDwedar/blogo/internal/usecase"
	"AbdelrahmanDwedar/blogo/pkg/auth"
	"AbdelrahmanDwedar/blogo/pkg/response"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userUC *usecase.UserUseCase
}

// NewUserHandler creates a new user handler
func NewUserHandler(userUC *usecase.UserUseCase) *UserHandler {
	return &UserHandler{userUC: userUC}
}

// CreateUser creates a new user
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username    string `json:"username"`
		Email       string `json:"email"`
		DisplayName string `json:"display_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.Username == "" || req.Email == "" || req.DisplayName == "" {
		response.Error(w, http.StatusBadRequest, "Username, email, and display_name are required")
		return
	}

	user, token, err := h.userUC.CreateUser(req.Username, req.Email, req.DisplayName)
	if err != nil {
		if err == entity.ErrInvalidUsername || err == entity.ErrInvalidEmail || err == entity.ErrInvalidDisplayName {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	response.Created(w, map[string]interface{}{
		"user":  user,
		"token": token,
	})
}

// GetUser retrieves a user by ID
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID, err := getIDFromPath(r, "id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, stats, err := h.userUC.GetUserWithStats(userID)
	if err != nil {
		if err == entity.ErrUserNotFound {
			response.Error(w, http.StatusNotFound, "User not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	response.Success(w, map[string]interface{}{
		"user":  user,
		"stats": stats,
	})
}

// UpdateUser updates user information
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	claims, err := auth.GetUserFromContext(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userID, err := getIDFromPath(r, "id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Check if user is updating their own profile
	if claims.UserID != userID {
		response.Error(w, http.StatusForbidden, "You can only update your own profile")
		return
	}

	var req struct {
		DisplayName  string `json:"display_name"`
		Bio          string `json:"bio"`
		ProfileImage string `json:"profile_image"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, err := h.userUC.UpdateUser(userID, req.DisplayName, req.Bio, req.ProfileImage)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to update user")
		return
	}

	response.Success(w, user)
}

// FollowUser follows or unfollows a user
func (h *UserHandler) FollowUser(w http.ResponseWriter, r *http.Request) {
	claims, err := auth.GetUserFromContext(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userID, err := getIDFromPath(r, "id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req struct {
		Action string `json:"action"` // "follow" or "unfollow"
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.Action == "follow" {
		err = h.userUC.FollowUser(claims.UserID, userID)
	} else if req.Action == "unfollow" {
		err = h.userUC.UnfollowUser(claims.UserID, userID)
	} else {
		response.Error(w, http.StatusBadRequest, "Invalid action. Use 'follow' or 'unfollow'")
		return
	}

	if err != nil {
		if err == entity.ErrCannotFollowSelf {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "Failed to "+req.Action+" user")
		return
	}

	response.Success(w, map[string]string{
		"message": "Successfully " + req.Action + "ed user",
	})
}

// GetUserFollowers retrieves a user's followers
func (h *UserHandler) GetUserFollowers(w http.ResponseWriter, r *http.Request) {
	userID, err := getIDFromPath(r, "id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	limit, offset := getPaginationParams(r)

	followers, err := h.userUC.GetFollowers(userID, limit, offset)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to get followers")
		return
	}

	response.Success(w, map[string]interface{}{
		"followers": followers,
		"limit":     limit,
		"offset":    offset,
	})
}

// GetUserFollowing retrieves users that a user is following
func (h *UserHandler) GetUserFollowing(w http.ResponseWriter, r *http.Request) {
	userID, err := getIDFromPath(r, "id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	limit, offset := getPaginationParams(r)

	following, err := h.userUC.GetFollowing(userID, limit, offset)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to get following")
		return
	}

	response.Success(w, map[string]interface{}{
		"following": following,
		"limit":     limit,
		"offset":    offset,
	})
}


