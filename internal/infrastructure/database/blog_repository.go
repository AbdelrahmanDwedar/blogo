package database

import (
	"database/sql"
	"fmt"

	"AbdelrahmanDwedar/blogo/internal/domain/entity"
)

// BlogRepository implements the blog repository interface
type BlogRepository struct {
	db *PostgresDB
}

// NewBlogRepository creates a new blog repository
func NewBlogRepository(db *PostgresDB) *BlogRepository {
	return &BlogRepository{db: db}
}

// Create creates a new blog post
func (r *BlogRepository) Create(blog *entity.Blog) error {
	r.db.mu.Lock()
	defer r.db.mu.Unlock()

	err := r.db.Client.QueryRow(`
		INSERT INTO blogs (title, description, body, author_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`, blog.Title, blog.Description, blog.Body, blog.AuthorID,
		blog.CreatedAt, blog.UpdatedAt).Scan(&blog.ID)

	if err != nil {
		return fmt.Errorf("create blog: %w", err)
	}
	return nil
}

// GetByID retrieves a blog by ID
func (r *BlogRepository) GetByID(id int64) (*entity.Blog, error) {
	r.db.mu.RLock()
	defer r.db.mu.RUnlock()

	blog := &entity.Blog{Author: &entity.User{}}
	err := r.db.Client.QueryRow(`
		SELECT b.id, b.title, b.description, b.body, b.author_id,
		       u.id, u.username, u.email, u.display_name, u.bio, u.profile_image, u.created_at, u.updated_at,
		       (SELECT COUNT(*) FROM likes WHERE blog_id = b.id) as likes_count,
		       b.created_at, b.updated_at
		FROM blogs b
		INNER JOIN users u ON b.author_id = u.id
		WHERE b.id = $1
	`, id).Scan(
		&blog.ID, &blog.Title, &blog.Description, &blog.Body, &blog.AuthorID,
		&blog.Author.ID, &blog.Author.Username, &blog.Author.Email, &blog.Author.DisplayName,
		&blog.Author.Bio, &blog.Author.ProfileImage, &blog.Author.CreatedAt, &blog.Author.UpdatedAt,
		&blog.LikesCount, &blog.CreatedAt, &blog.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, entity.ErrBlogNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get blog by id: %w", err)
	}

	return blog, nil
}

// GetAll retrieves all blogs with pagination
func (r *BlogRepository) GetAll(limit, offset int) ([]*entity.Blog, error) {
	r.db.mu.RLock()
	defer r.db.mu.RUnlock()

	rows, err := r.db.Client.Query(`
		SELECT b.id, b.title, b.description, b.body, b.author_id,
		       u.id, u.username, u.email, u.display_name, u.bio, u.profile_image, u.created_at, u.updated_at,
		       (SELECT COUNT(*) FROM likes WHERE blog_id = b.id) as likes_count,
		       b.created_at, b.updated_at
		FROM blogs b
		INNER JOIN users u ON b.author_id = u.id
		ORDER BY b.created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get all blogs: %w", err)
	}
	defer rows.Close()

	blogs := []*entity.Blog{}
	for rows.Next() {
		blog := &entity.Blog{Author: &entity.User{}}
		err := rows.Scan(
			&blog.ID, &blog.Title, &blog.Description, &blog.Body, &blog.AuthorID,
			&blog.Author.ID, &blog.Author.Username, &blog.Author.Email, &blog.Author.DisplayName,
			&blog.Author.Bio, &blog.Author.ProfileImage, &blog.Author.CreatedAt, &blog.Author.UpdatedAt,
			&blog.LikesCount, &blog.CreatedAt, &blog.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan blog: %w", err)
		}
		blogs = append(blogs, blog)
	}

	return blogs, nil
}

// GetByAuthor retrieves blogs by a specific author
func (r *BlogRepository) GetByAuthor(authorID int64, limit, offset int) ([]*entity.Blog, error) {
	r.db.mu.RLock()
	defer r.db.mu.RUnlock()

	rows, err := r.db.Client.Query(`
		SELECT b.id, b.title, b.description, b.body, b.author_id,
		       u.id, u.username, u.email, u.display_name, u.bio, u.profile_image, u.created_at, u.updated_at,
		       (SELECT COUNT(*) FROM likes WHERE blog_id = b.id) as likes_count,
		       b.created_at, b.updated_at
		FROM blogs b
		INNER JOIN users u ON b.author_id = u.id
		WHERE b.author_id = $1
		ORDER BY b.created_at DESC
		LIMIT $2 OFFSET $3
	`, authorID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get blogs by author: %w", err)
	}
	defer rows.Close()

	blogs := []*entity.Blog{}
	for rows.Next() {
		blog := &entity.Blog{Author: &entity.User{}}
		err := rows.Scan(
			&blog.ID, &blog.Title, &blog.Description, &blog.Body, &blog.AuthorID,
			&blog.Author.ID, &blog.Author.Username, &blog.Author.Email, &blog.Author.DisplayName,
			&blog.Author.Bio, &blog.Author.ProfileImage, &blog.Author.CreatedAt, &blog.Author.UpdatedAt,
			&blog.LikesCount, &blog.CreatedAt, &blog.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan blog: %w", err)
		}
		blogs = append(blogs, blog)
	}

	return blogs, nil
}

// Update updates a blog post
func (r *BlogRepository) Update(blog *entity.Blog) error {
	r.db.mu.Lock()
	defer r.db.mu.Unlock()

	_, err := r.db.Client.Exec(`
		UPDATE blogs
		SET title = $1, description = $2, body = $3, updated_at = $4
		WHERE id = $5 AND author_id = $6
	`, blog.Title, blog.Description, blog.Body, blog.UpdatedAt, blog.ID, blog.AuthorID)

	if err != nil {
		return fmt.Errorf("update blog: %w", err)
	}
	return nil
}

// Delete deletes a blog post
func (r *BlogRepository) Delete(id, authorID int64) error {
	r.db.mu.Lock()
	defer r.db.mu.Unlock()

	result, err := r.db.Client.Exec(`
		DELETE FROM blogs
		WHERE id = $1 AND author_id = $2
	`, id, authorID)
	if err != nil {
		return fmt.Errorf("delete blog: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return entity.ErrBlogNotFound
	}

	return nil
}

// Like adds a like to a blog
func (r *BlogRepository) Like(blogID, userID int64) error {
	r.db.mu.Lock()
	defer r.db.mu.Unlock()

	_, err := r.db.Client.Exec(`
		INSERT INTO likes (blog_id, user_id, created_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (blog_id, user_id) DO NOTHING
	`, blogID, userID)

	if err != nil {
		return fmt.Errorf("like blog: %w", err)
	}
	return nil
}

// Unlike removes a like from a blog
func (r *BlogRepository) Unlike(blogID, userID int64) error {
	r.db.mu.Lock()
	defer r.db.mu.Unlock()

	_, err := r.db.Client.Exec(`
		DELETE FROM likes
		WHERE blog_id = $1 AND user_id = $2
	`, blogID, userID)

	if err != nil {
		return fmt.Errorf("unlike blog: %w", err)
	}
	return nil
}

// GetLikes retrieves users who liked a blog
func (r *BlogRepository) GetLikes(blogID int64, limit, offset int) ([]*entity.User, error) {
	r.db.mu.RLock()
	defer r.db.mu.RUnlock()

	rows, err := r.db.Client.Query(`
		SELECT u.id, u.username, u.email, u.display_name, u.bio, u.profile_image, u.created_at, u.updated_at
		FROM users u
		INNER JOIN likes l ON u.id = l.user_id
		WHERE l.blog_id = $1
		ORDER BY l.created_at DESC
		LIMIT $2 OFFSET $3
	`, blogID, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("get blog likes: %w", err)
	}
	defer rows.Close()

	users := []*entity.User{}
	for rows.Next() {
		user := &entity.User{}
		err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.DisplayName,
			&user.Bio, &user.ProfileImage, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}

// IsLikedBy checks if a user has liked a blog
func (r *BlogRepository) IsLikedBy(blogID, userID int64) (bool, error) {
	r.db.mu.RLock()
	defer r.db.mu.RUnlock()

	var exists bool
	err := r.db.Client.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM likes
			WHERE blog_id = $1 AND user_id = $2
		)
	`, blogID, userID).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("check like: %w", err)
	}
	return exists, nil
}


