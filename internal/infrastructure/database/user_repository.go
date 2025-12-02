package database

import (
	"database/sql"
	"fmt"

	"AbdelrahmanDwedar/blogo/internal/domain/entity"
)

// UserRepository implements the user repository interface
type UserRepository struct {
	db *PostgresDB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *PostgresDB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(user *entity.User) error {
	r.db.mu.Lock()
	defer r.db.mu.Unlock()

	err := r.db.Client.QueryRow(`
		INSERT INTO users (username, email, display_name, bio, profile_image, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`, user.Username, user.Email, user.DisplayName, user.Bio, user.ProfileImage,
		user.CreatedAt, user.UpdatedAt).Scan(&user.ID)

	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(id int64) (*entity.User, error) {
	r.db.mu.RLock()
	defer r.db.mu.RUnlock()

	user := &entity.User{}
	err := r.db.Client.QueryRow(`
		SELECT id, username, email, display_name, bio, profile_image, created_at, updated_at
		FROM users
		WHERE id = $1
	`, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.DisplayName,
		&user.Bio, &user.ProfileImage, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, entity.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	return user, nil
}

// GetByUsername retrieves a user by username
func (r *UserRepository) GetByUsername(username string) (*entity.User, error) {
	r.db.mu.RLock()
	defer r.db.mu.RUnlock()

	user := &entity.User{}
	err := r.db.Client.QueryRow(`
		SELECT id, username, email, display_name, bio, profile_image, created_at, updated_at
		FROM users
		WHERE username = $1
	`, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.DisplayName,
		&user.Bio, &user.ProfileImage, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, entity.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get user by username: %w", err)
	}

	return user, nil
}

// Update updates user information
func (r *UserRepository) Update(user *entity.User) error {
	r.db.mu.Lock()
	defer r.db.mu.Unlock()

	_, err := r.db.Client.Exec(`
		UPDATE users
		SET display_name = $1, bio = $2, profile_image = $3, updated_at = $4
		WHERE id = $5
	`, user.DisplayName, user.Bio, user.ProfileImage, user.UpdatedAt, user.ID)

	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}

// Follow creates a follow relationship
func (r *UserRepository) Follow(followerID, followingID int64) error {
	r.db.mu.Lock()
	defer r.db.mu.Unlock()

	_, err := r.db.Client.Exec(`
		INSERT INTO followers (follower_id, following_id, created_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (follower_id, following_id) DO NOTHING
	`, followerID, followingID)

	if err != nil {
		return fmt.Errorf("follow user: %w", err)
	}
	return nil
}

// Unfollow removes a follow relationship
func (r *UserRepository) Unfollow(followerID, followingID int64) error {
	r.db.mu.Lock()
	defer r.db.mu.Unlock()

	_, err := r.db.Client.Exec(`
		DELETE FROM followers
		WHERE follower_id = $1 AND following_id = $2
	`, followerID, followingID)

	if err != nil {
		return fmt.Errorf("unfollow user: %w", err)
	}
	return nil
}

// GetFollowers retrieves users following the given user
func (r *UserRepository) GetFollowers(userID int64, limit, offset int) ([]*entity.User, error) {
	r.db.mu.RLock()
	defer r.db.mu.RUnlock()

	rows, err := r.db.Client.Query(`
		SELECT u.id, u.username, u.email, u.display_name, u.bio, u.profile_image, u.created_at, u.updated_at
		FROM users u
		INNER JOIN followers f ON u.id = f.follower_id
		WHERE f.following_id = $1
		ORDER BY f.created_at DESC
		LIMIT $2 OFFSET $3
	`, userID, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("get followers: %w", err)
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
			return nil, fmt.Errorf("scan follower: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}

// GetFollowing retrieves users that the given user is following
func (r *UserRepository) GetFollowing(userID int64, limit, offset int) ([]*entity.User, error) {
	r.db.mu.RLock()
	defer r.db.mu.RUnlock()

	rows, err := r.db.Client.Query(`
		SELECT u.id, u.username, u.email, u.display_name, u.bio, u.profile_image, u.created_at, u.updated_at
		FROM users u
		INNER JOIN followers f ON u.id = f.following_id
		WHERE f.follower_id = $1
		ORDER BY f.created_at DESC
		LIMIT $2 OFFSET $3
	`, userID, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("get following: %w", err)
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
			return nil, fmt.Errorf("scan following: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}

// IsFollowing checks if one user is following another
func (r *UserRepository) IsFollowing(followerID, followingID int64) (bool, error) {
	r.db.mu.RLock()
	defer r.db.mu.RUnlock()

	var exists bool
	err := r.db.Client.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM followers
			WHERE follower_id = $1 AND following_id = $2
		)
	`, followerID, followingID).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("check following: %w", err)
	}
	return exists, nil
}

// GetStats retrieves statistics for a user
func (r *UserRepository) GetStats(userID int64) (*entity.UserStats, error) {
	r.db.mu.RLock()
	defer r.db.mu.RUnlock()

	stats := &entity.UserStats{}

	// Get followers count
	err := r.db.Client.QueryRow(`
		SELECT COUNT(*) FROM followers WHERE following_id = $1
	`, userID).Scan(&stats.FollowersCount)
	if err != nil {
		return nil, fmt.Errorf("get followers count: %w", err)
	}

	// Get following count
	err = r.db.Client.QueryRow(`
		SELECT COUNT(*) FROM followers WHERE follower_id = $1
	`, userID).Scan(&stats.FollowingCount)
	if err != nil {
		return nil, fmt.Errorf("get following count: %w", err)
	}

	// Get blogs count
	err = r.db.Client.QueryRow(`
		SELECT COUNT(*) FROM blogs WHERE author_id = $1
	`, userID).Scan(&stats.BlogsCount)
	if err != nil {
		return nil, fmt.Errorf("get blogs count: %w", err)
	}

	return stats, nil
}


