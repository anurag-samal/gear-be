package auth

import (
	"github.com/google/uuid"
	"time"
)

type BaseModel struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

type User struct {
	BaseModel
	OrganizationID  *uuid.UUID
	Email           string
	FullName        string
	AvatarURL       *string
	IsEmailVerified bool
	LastLoginAt     *time.Time
}

type UserCredential struct {
	UserID       uuid.UUID
	PasswordHash string
	CreatedAt    time.Time
}

type AuthAccount struct {
	BaseModel
	UserID         uuid.UUID
	Provider       string
	ProviderUserID string
}

type RefreshToken struct {
	BaseModel
	UserID    uuid.UUID
	TokenHash string
	ExpiresAt time.Time
	Revoked   bool
}

type Organization struct {
	BaseModel
	Name         string
	BusinessType string
	Industry     *string
	Country      *string
	OwnerUserID  *uuid.UUID
}

type Team struct {
	BaseModel
	OrganizationID uuid.UUID
	Name           string
	Description    *string
}

type TeamMember struct {
	TeamID   uuid.UUID
	UserID   uuid.UUID
	Role     string
	JoinedAt time.Time
}

type Invitation struct {
	BaseModel
	OrganizationID uuid.UUID
	TeamID         uuid.UUID
	Email          string
	InvitedBy      uuid.UUID
	Status         string
	ExpiresAt      time.Time
}

type Subscription struct {
	BaseModel
	OrganizationID     uuid.UUID
	Plan               string
	Status             string
	TrialEndsAt        *time.Time
	CurrentPeriodStart *time.Time
	CurrentPeriodEnd   *time.Time
}
