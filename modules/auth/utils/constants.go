package auth

import (
	"errors"
	"time"
)

const (
	JWT_ISSUER            = "gear"
	JWT_ACCESS_TOKEN_TTL  = 15 * time.Minute
	JWT_REFRESH_TOKEN_TTL = 7 * 24 * time.Hour
	PKCE_TTL              = 10 * time.Minute
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidState       = errors.New("invalid state")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")

	ErrRefreshTokenNotFound = errors.New("refresh token not found")
	ErrRefreshTokenInvalid  = errors.New("invalid refresh token")
	ErrRefreshTokenExpired  = errors.New("refresh token expired")
	ErrRefreshTokenRevoked  = errors.New("refresh token revoked")

	ErrUnauthorized = errors.New("user is unauthorized")
)
