package repository

import (
	"time"

	"AbdelrahmanDwedar/blogo/internal/domain/entity"
)

// CacheRepository defines the interface for caching operations
type CacheRepository interface {
	// User cache operations
	SetUser(user *entity.User, expiration time.Duration) error
	GetUser(userID int64) (*entity.User, error)
	DeleteUser(userID int64) error

	// Blog cache operations
	SetBlog(blog *entity.Blog, expiration time.Duration) error
	GetBlog(blogID int64) (*entity.Blog, error)
	DeleteBlog(blogID int64) error

	// Bulk operations
	DeletePattern(pattern string) error
}


