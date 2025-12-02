package entity

import "time"

// Blog represents a blog post entity in the domain
type Blog struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Body        string    `json:"body"`
	AuthorID    int64     `json:"author_id"`
	Author      *User     `json:"author,omitempty"`
	LikesCount  int       `json:"likes_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewBlog creates a new blog entity
func NewBlog(title, description, body string, authorID int64) *Blog {
	now := time.Now()
	return &Blog{
		Title:       title,
		Description: description,
		Body:        body,
		AuthorID:    authorID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// Update updates blog information
func (b *Blog) Update(title, description, body string) {
	b.Title = title
	b.Description = description
	b.Body = body
	b.UpdatedAt = time.Now()
}

// Validate validates blog data
func (b *Blog) Validate() error {
	if b.Title == "" {
		return ErrInvalidTitle
	}
	if b.Body == "" {
		return ErrInvalidBody
	}
	if b.AuthorID == 0 {
		return ErrInvalidAuthor
	}
	return nil
}

// IsOwnedBy checks if the blog is owned by the given user
func (b *Blog) IsOwnedBy(userID int64) bool {
	return b.AuthorID == userID
}


