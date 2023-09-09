package tables

import (
	"fmt"
	"time"
)

type BlogStorer interface {
	GetBlogs(b Blog) (*[]Blog, error)
	GetBlogById(id int64) (*Blog, error)
	CreateBlog() (*Blog, error)
	EditeBlog(id int64, b Blog) (*Blog, error)
	DeleteBlog(id int64) (*Blog, error)
}

type Blog struct {
	ID          int64       `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Body        string      `json:"body"`
	PublishedAt time.Time   `json:"published_at"`
	Auther      interface{} `json:"auther"`
}

func (s *PostgresStore) GetBlogs() (*[]Blog, error) {
	s.ru.RLock()
	defer s.ru.RUnlock()

	row, err := s.Client.Query(`
		SELECT * FROM blog
	`)
	if err != nil {
		return nil, fmt.Errorf("Blog: %s", err)
	}
	defer row.Close()

	blogs := []Blog{}

	for row.Next() {
		blog := Blog{}
		err := row.Scan(&blog.ID, &blog.Title, &blog.Description, &blog.Body, &blog.PublishedAt, &blog.Auther)
		if err != nil {
			return nil, fmt.Errorf("Blog: %s", err)
		}
		blogs = append(blogs, blog)
	}

	return &blogs, nil
}

func (s *PostgresStore) GetBlogById(id int64) (*Blog, error) {
	s.ru.RLock()
	defer s.ru.RUnlock()

	row, err := s.Client.Query(`
		SELECT * FROM blog
		WHERE ?
	`, id)
	if err != nil {
		return nil, fmt.Errorf("Blog: %s", err)
	}
	defer row.Close()

	blog := Blog{}

	err = row.Scan(&blog.ID, &blog.Title, &blog.Description, &blog.Body, &blog.PublishedAt, &blog.Auther)
	if err != nil {
		return nil, fmt.Errorf("Blog: %s", err)
	}

	return &blog, nil
}

func (s *PostgresStore) CreateBlog(b Blog) (*Blog, error) {
	s.ru.Lock()
	defer s.ru.Unlock()

	stmt, err := s.Client.Prepare(`
		INSERT INTO blog
		(title, Description, body, auther)
		VALUES (?, ?, ?, ?)
	`)
	if err != nil {
		return nil, fmt.Errorf("Blog: %s", err)
	}

	res, err := stmt.Exec(b.Title, b.Description, b.Body, b.Auther)
	if err != nil {
		return nil, fmt.Errorf("Blog: %s", err)
	}
	defer stmt.Close()

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("Blog: %s", err)
	}

	return &Blog{
		ID:          id,
		Title:       b.Title,
		Description: b.Description,
		Body:        b.Body,
		PublishedAt: time.Now(),
		Auther:      b.Auther,
	}, nil
}

func (s *PostgresStore) EditeBlog(id int64, b Blog) (*Blog, error) {
	s.ru.Lock()
	defer s.ru.Unlock()

	stmt, err := s.Client.Prepare(`
		UPDATE blog
		SET title = ?, description = ?, body = ?, auther = ?
		WHERE id = ?
	`)
	if err != nil {
		return nil, fmt.Errorf("Blog: %s", err)
	}

	_, err = stmt.Exec(b.Title, b.Description, b.Body, b.Auther, id)
	if err != nil {
		return nil, fmt.Errorf("Blog: %s", err)
	}
	defer stmt.Close()

	return &Blog{
		ID:          id,
		Title:       b.Title,
		Description: b.Description,
		Body:        b.Body,
		Auther:      b.Auther,
	}, nil
}

func (s *PostgresStore) DeleteBlog(id int64) error {
	s.ru.Lock()
	defer s.ru.Unlock()

	stmt, err := s.Client.Prepare(`
		DELETE FROM blog
		WHERE id = ?
	`)
	if err != nil {
		return fmt.Errorf("Blog: %s", err)
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("Blog: %s", err)
	}
	defer stmt.Close()

	return nil
}
