package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	_ "github.com/lib/pq"
)

// PostgresDB represents a PostgreSQL database connection
type PostgresDB struct {
	Client *sql.DB
	mu     *sync.RWMutex
}

// NewPostgresDB creates a new PostgreSQL database connection
func NewPostgresDB() (*PostgresDB, error) {
	connStr := fmt.Sprintf(
		"user=%s dbname=%s password=%s host=%s port=%s sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("✅ PostgreSQL connected successfully")

	return &PostgresDB{
		Client: db,
		mu:     &sync.RWMutex{},
	}, nil
}

// Close closes the database connection
func (db *PostgresDB) Close() error {
	if db.Client != nil {
		return db.Client.Close()
	}
	return nil
}

// InitTables creates all necessary database tables
func (db *PostgresDB) InitTables() error {
	// Create users table
	_, err := db.Client.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) UNIQUE NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			display_name VARCHAR(100) NOT NULL,
			bio TEXT,
			profile_image VARCHAR(500),
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("create users table: %w", err)
	}

	// Create blogs table
	_, err = db.Client.Exec(`
		CREATE TABLE IF NOT EXISTS blogs (
			id SERIAL PRIMARY KEY,
			title VARCHAR(200) NOT NULL,
			description TEXT,
			body TEXT NOT NULL,
			author_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("create blogs table: %w", err)
	}

	// Create followers table
	_, err = db.Client.Exec(`
		CREATE TABLE IF NOT EXISTS followers (
			id SERIAL PRIMARY KEY,
			follower_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			following_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(follower_id, following_id),
			CHECK (follower_id != following_id)
		)
	`)
	if err != nil {
		return fmt.Errorf("create followers table: %w", err)
	}

	// Create likes table
	_, err = db.Client.Exec(`
		CREATE TABLE IF NOT EXISTS likes (
			id SERIAL PRIMARY KEY,
			blog_id INTEGER NOT NULL REFERENCES blogs(id) ON DELETE CASCADE,
			user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(blog_id, user_id)
		)
	`)
	if err != nil {
		return fmt.Errorf("create likes table: %w", err)
	}

	// Create indexes for better performance
	_, err = db.Client.Exec(`
		CREATE INDEX IF NOT EXISTS idx_blogs_author ON blogs(author_id);
		CREATE INDEX IF NOT EXISTS idx_blogs_created ON blogs(created_at DESC);
		CREATE INDEX IF NOT EXISTS idx_followers_follower ON followers(follower_id);
		CREATE INDEX IF NOT EXISTS idx_followers_following ON followers(following_id);
		CREATE INDEX IF NOT EXISTS idx_likes_blog ON likes(blog_id);
		CREATE INDEX IF NOT EXISTS idx_likes_user ON likes(user_id);
	`)
	if err != nil {
		return fmt.Errorf("create indexes: %w", err)
	}

	log.Println("✅ Database tables initialized successfully")
	return nil
}


