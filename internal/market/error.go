package market

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrOrderIsExist       = errors.New("order is loaded by other User")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAccessToken = errors.New("invalid access token")
	ErrFormatOrderID      = errors.New("invalid access token")
)