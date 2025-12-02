package entity

import "errors"

// Domain errors
var (
	// User errors
	ErrInvalidUsername    = errors.New("invalid username")
	ErrInvalidEmail       = errors.New("invalid email")
	ErrInvalidDisplayName = errors.New("invalid display name")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrCannotFollowSelf   = errors.New("cannot follow yourself")
	ErrUnauthorized       = errors.New("unauthorized")

	// Blog errors
	ErrInvalidTitle  = errors.New("invalid title")
	ErrInvalidBody   = errors.New("invalid body")
	ErrInvalidAuthor = errors.New("invalid author")
	ErrBlogNotFound  = errors.New("blog not found")
	ErrNotBlogOwner  = errors.New("not blog owner")

	// General errors
	ErrInvalidID = errors.New("invalid ID")
)


