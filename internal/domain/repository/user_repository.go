package repository

import "AbdelrahmanDwedar/blogo/internal/domain/entity"

// UserRepository defines the interface for user data access
type UserRepository interface {
	// Create creates a new user
	Create(user *entity.User) error

	// GetByID retrieves a user by ID
	GetByID(id int64) (*entity.User, error)

	// GetByUsername retrieves a user by username
	GetByUsername(username string) (*entity.User, error)

	// Update updates user information
	Update(user *entity.User) error

	// Follow creates a follow relationship
	Follow(followerID, followingID int64) error

	// Unfollow removes a follow relationship
	Unfollow(followerID, followingID int64) error

	// GetFollowers retrieves users following the given user
	GetFollowers(userID int64, limit, offset int) ([]*entity.User, error)

	// GetFollowing retrieves users that the given user is following
	GetFollowing(userID int64, limit, offset int) ([]*entity.User, error)

	// IsFollowing checks if one user is following another
	IsFollowing(followerID, followingID int64) (bool, error)

	// GetStats retrieves statistics for a user
	GetStats(userID int64) (*entity.UserStats, error)
}


