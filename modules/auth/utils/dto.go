package auth

import "github.com/google/uuid"

type RegisterRequest struct {
	FullName string `json:"full_name" validate:"required,min=3,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required"`
}

type GoogleLoginRequest struct {
	Code     string `json:"code" validate:"required"`
	State    string `json:"state" validate:"required"`
	Verifier string `json:"verifier" validate:"required"`
}

type AuthResponse struct {
	AccessToken string          `json:"access_token"`
	TokenType   string          `json:"token_type"`
	ExpiresIn   int64           `json:"expires_in"`
	User        ProfileResponse `json:"user"`
}

type ProfileResponse struct {
	ID              uuid.UUID  `json:"id"`
	OrganizationID  *uuid.UUID `json:"organization_id,omitempty"`
	Email           string     `json:"email"`
	FullName        string     `json:"full_name"`
	AvatarURL       *string    `json:"avatar_url,omitempty"`
	IsEmailVerified bool       `json:"is_email_verified"`
}
