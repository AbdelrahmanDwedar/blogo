package entity

import "time"

// User represents a user entity in the domain
type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	DisplayName  string    `json:"display_name"`
	Bio          string    `json:"bio,omitempty"`
	ProfileImage string    `json:"profile_image,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// UserStats represents user statistics
type UserStats struct {
	FollowersCount int `json:"followers_count"`
	FollowingCount int `json:"following_count"`
	BlogsCount     int `json:"blogs_count"`
}

// NewUser creates a new user entity
func NewUser(username, email, displayName string) *User {
	now := time.Now()
	return &User{
		Username:    username,
		Email:       email,
		DisplayName: displayName,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// Update updates user information
func (u *User) Update(displayName, bio, profileImage string) {
	u.DisplayName = displayName
	u.Bio = bio
	u.ProfileImage = profileImage
	u.UpdatedAt = time.Now()
}

// Validate validates user data
func (u *User) Validate() error {
	if u.Username == "" {
		return ErrInvalidUsername
	}
	if u.Email == "" {
		return ErrInvalidEmail
	}
	if u.DisplayName == "" {
		return ErrInvalidDisplayName
	}
	return nil
}


