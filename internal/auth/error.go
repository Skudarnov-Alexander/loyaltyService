package auth

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserIsExist        = errors.New("user is exist")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAccessToken = errors.New("invalid access token")
)