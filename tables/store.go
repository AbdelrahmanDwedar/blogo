package tables

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
)

type PostgresStore struct {
	Client   *sql.DB
	ru       *sync.RWMutex
	user     string
	name     string
	password string
	host     string
	port     string
}

func NewPostgresStore() *PostgresStore {
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
		log.Fatal(err)
	}

	return &PostgresStore{
		Client:   db,
		user:     os.Getenv("DB_USER"),
		name:     os.Getenv("DB_NAME"),
		password: os.Getenv("DB_PASSWORD"),
		host:     os.Getenv("DB_HOST"),
		port:     os.Getenv("DB_PORT"),
	}
}

func (s *PostgresStore) CreateBlogsTables() error {
	stmt, err := s.Client.Prepare(`
		CTEARE TABLE blog (
			id int NOT NULL PRIMARY KEY,
			title varchar(20) NOT NULL,
			description text,
			body longtext NOT NULL,
			published_at DATETIME,
			auther int FOREIGN KEY REFERENCES user(id)
		)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	return nil
}
