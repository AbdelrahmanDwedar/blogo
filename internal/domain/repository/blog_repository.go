package repository

import "AbdelrahmanDwedar/blogo/internal/domain/entity"

// BlogRepository defines the interface for blog data access
type BlogRepository interface {
	// Create creates a new blog post
	Create(blog *entity.Blog) error

	// GetByID retrieves a blog by ID
	GetByID(id int64) (*entity.Blog, error)

	// GetAll retrieves all blogs with pagination
	GetAll(limit, offset int) ([]*entity.Blog, error)

	// GetByAuthor retrieves blogs by a specific author
	GetByAuthor(authorID int64, limit, offset int) ([]*entity.Blog, error)

	// Update updates a blog post
	Update(blog *entity.Blog) error

	// Delete deletes a blog post
	Delete(id, authorID int64) error

	// Like adds a like to a blog
	Like(blogID, userID int64) error

	// Unlike removes a like from a blog
	Unlike(blogID, userID int64) error

	// GetLikes retrieves users who liked a blog
	GetLikes(blogID int64, limit, offset int) ([]*entity.User, error)

	// IsLikedBy checks if a user has liked a blog
	IsLikedBy(blogID, userID int64) (bool, error)
}


