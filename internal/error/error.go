package error

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserIsExist  = errors.New("user is exist")

	ErrOrderIsExist  = errors.New("order is loaded by other User")
	ErrFormatOrderID = errors.New("order format is incorrect")

	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAccessToken = errors.New("invalid access token")
	ErrClaimsParsing      = errors.New("couldn't parse claims")
	ErrExpiredToken       = errors.New("token is expired")
)
