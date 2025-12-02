package usecase

import (
	"time"

	"AbdelrahmanDwedar/blogo/internal/domain/entity"
	"AbdelrahmanDwedar/blogo/internal/domain/repository"
	"AbdelrahmanDwedar/blogo/pkg/auth"
)

// UserUseCase handles user-related business logic
type UserUseCase struct {
	userRepo  repository.UserRepository
	cacheRepo repository.CacheRepository
}

// NewUserUseCase creates a new user use case
func NewUserUseCase(userRepo repository.UserRepository, cacheRepo repository.CacheRepository) *UserUseCase {
	return &UserUseCase{
		userRepo:  userRepo,
		cacheRepo: cacheRepo,
	}
}

// CreateUser creates a new user and returns JWT token
func (uc *UserUseCase) CreateUser(username, email, displayName string) (*entity.User, string, error) {
	// Create user entity
	user := entity.NewUser(username, email, displayName)

	// Validate
	if err := user.Validate(); err != nil {
		return nil, "", err
	}

	// Save to database
	if err := uc.userRepo.Create(user); err != nil {
		return nil, "", err
	}

	// Cache the user
	if uc.cacheRepo != nil {
		uc.cacheRepo.SetUser(user, 15*time.Minute)
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, user.Username, user.Email)
	if err != nil {
		return user, "", err
	}

	return user, token, nil
}

// GetUserByID retrieves a user by ID with caching
func (uc *UserUseCase) GetUserByID(id int64) (*entity.User, error) {
	// Try cache first
	if uc.cacheRepo != nil {
		if user, err := uc.cacheRepo.GetUser(id); err == nil && user != nil {
			return user, nil
		}
	}

	// Get from database
	user, err := uc.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Cache it
	if uc.cacheRepo != nil {
		uc.cacheRepo.SetUser(user, 15*time.Minute)
	}

	return user, nil
}

// GetUserWithStats retrieves a user with statistics
func (uc *UserUseCase) GetUserWithStats(id int64) (*entity.User, *entity.UserStats, error) {
	user, err := uc.GetUserByID(id)
	if err != nil {
		return nil, nil, err
	}

	stats, err := uc.userRepo.GetStats(id)
	if err != nil {
		return user, nil, err
	}

	return user, stats, nil
}

// UpdateUser updates user profile
func (uc *UserUseCase) UpdateUser(id int64, displayName, bio, profileImage string) (*entity.User, error) {
	// Get existing user
	user, err := uc.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update
	user.Update(displayName, bio, profileImage)

	// Save
	if err := uc.userRepo.Update(user); err != nil {
		return nil, err
	}

	// Invalidate cache
	if uc.cacheRepo != nil {
		uc.cacheRepo.DeleteUser(id)
	}

	return user, nil
}

// FollowUser creates a follow relationship
func (uc *UserUseCase) FollowUser(followerID, followingID int64) error {
	if followerID == followingID {
		return entity.ErrCannotFollowSelf
	}

	if err := uc.userRepo.Follow(followerID, followingID); err != nil {
		return err
	}

	// Invalidate caches
	if uc.cacheRepo != nil {
		uc.cacheRepo.DeleteUser(followerID)
		uc.cacheRepo.DeleteUser(followingID)
	}

	return nil
}

// UnfollowUser removes a follow relationship
func (uc *UserUseCase) UnfollowUser(followerID, followingID int64) error {
	if err := uc.userRepo.Unfollow(followerID, followingID); err != nil {
		return err
	}

	// Invalidate caches
	if uc.cacheRepo != nil {
		uc.cacheRepo.DeleteUser(followerID)
		uc.cacheRepo.DeleteUser(followingID)
	}

	return nil
}

// GetFollowers retrieves a user's followers
func (uc *UserUseCase) GetFollowers(userID int64, limit, offset int) ([]*entity.User, error) {
	return uc.userRepo.GetFollowers(userID, limit, offset)
}

// GetFollowing retrieves users that a user is following
func (uc *UserUseCase) GetFollowing(userID int64, limit, offset int) ([]*entity.User, error) {
	return uc.userRepo.GetFollowing(userID, limit, offset)
}


