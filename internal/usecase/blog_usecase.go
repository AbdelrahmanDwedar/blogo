package usecase

import (
	"time"

	"AbdelrahmanDwedar/blogo/internal/domain/entity"
	"AbdelrahmanDwedar/blogo/internal/domain/repository"
)

// BlogUseCase handles blog-related business logic
type BlogUseCase struct {
	blogRepo  repository.BlogRepository
	cacheRepo repository.CacheRepository
}

// NewBlogUseCase creates a new blog use case
func NewBlogUseCase(blogRepo repository.BlogRepository, cacheRepo repository.CacheRepository) *BlogUseCase {
	return &BlogUseCase{
		blogRepo:  blogRepo,
		cacheRepo: cacheRepo,
	}
}

// CreateBlog creates a new blog post
func (uc *BlogUseCase) CreateBlog(title, description, body string, authorID int64) (*entity.Blog, error) {
	// Create blog entity
	blog := entity.NewBlog(title, description, body, authorID)

	// Validate
	if err := blog.Validate(); err != nil {
		return nil, err
	}

	// Save to database
	if err := uc.blogRepo.Create(blog); err != nil {
		return nil, err
	}

	// Invalidate blog list cache
	if uc.cacheRepo != nil {
		uc.cacheRepo.DeletePattern("blogs:*")
	}

	return blog, nil
}

// GetBlogByID retrieves a blog by ID with caching
func (uc *BlogUseCase) GetBlogByID(id int64) (*entity.Blog, error) {
	// Try cache first
	if uc.cacheRepo != nil {
		if blog, err := uc.cacheRepo.GetBlog(id); err == nil && blog != nil {
			return blog, nil
		}
	}

	// Get from database
	blog, err := uc.blogRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Cache it
	if uc.cacheRepo != nil {
		uc.cacheRepo.SetBlog(blog, 10*time.Minute)
	}

	return blog, nil
}

// GetAllBlogs retrieves all blogs with pagination
func (uc *BlogUseCase) GetAllBlogs(limit, offset int) ([]*entity.Blog, error) {
	return uc.blogRepo.GetAll(limit, offset)
}

// GetBlogsByAuthor retrieves blogs by a specific author
func (uc *BlogUseCase) GetBlogsByAuthor(authorID int64, limit, offset int) ([]*entity.Blog, error) {
	return uc.blogRepo.GetByAuthor(authorID, limit, offset)
}

// UpdateBlog updates a blog post
func (uc *BlogUseCase) UpdateBlog(id int64, title, description, body string, userID int64) (*entity.Blog, error) {
	// Get existing blog
	blog, err := uc.blogRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Check ownership
	if !blog.IsOwnedBy(userID) {
		return nil, entity.ErrNotBlogOwner
	}

	// Update
	blog.Update(title, description, body)

	// Validate
	if err := blog.Validate(); err != nil {
		return nil, err
	}

	// Save
	if err := uc.blogRepo.Update(blog); err != nil {
		return nil, err
	}

	// Invalidate caches
	if uc.cacheRepo != nil {
		uc.cacheRepo.DeleteBlog(id)
		uc.cacheRepo.DeletePattern("blogs:*")
	}

	return blog, nil
}

// DeleteBlog deletes a blog post
func (uc *BlogUseCase) DeleteBlog(id, userID int64) error {
	// Get blog to check ownership
	blog, err := uc.blogRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Check ownership
	if !blog.IsOwnedBy(userID) {
		return entity.ErrNotBlogOwner
	}

	// Delete
	if err := uc.blogRepo.Delete(id, userID); err != nil {
		return err
	}

	// Invalidate caches
	if uc.cacheRepo != nil {
		uc.cacheRepo.DeleteBlog(id)
		uc.cacheRepo.DeletePattern("blogs:*")
	}

	return nil
}

// LikeBlog adds a like to a blog
func (uc *BlogUseCase) LikeBlog(blogID, userID int64) error {
	if err := uc.blogRepo.Like(blogID, userID); err != nil {
		return err
	}

	// Invalidate cache
	if uc.cacheRepo != nil {
		uc.cacheRepo.DeleteBlog(blogID)
	}

	return nil
}

// UnlikeBlog removes a like from a blog
func (uc *BlogUseCase) UnlikeBlog(blogID, userID int64) error {
	if err := uc.blogRepo.Unlike(blogID, userID); err != nil {
		return err
	}

	// Invalidate cache
	if uc.cacheRepo != nil {
		uc.cacheRepo.DeleteBlog(blogID)
	}

	return nil
}

// GetBlogLikes retrieves users who liked a blog
func (uc *BlogUseCase) GetBlogLikes(blogID int64, limit, offset int) ([]*entity.User, error) {
	return uc.blogRepo.GetLikes(blogID, limit, offset)
}


