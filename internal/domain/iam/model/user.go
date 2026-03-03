package model

import (
	"time"

	"github.com/google/uuid"
)

// UserStatus represents the lifecycle state of a user account.
type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusInactive  UserStatus = "inactive"
	UserStatusSuspended UserStatus = "suspended"
)

// UserRole represents a user's role for access control.
type UserRole string

const (
	UserRoleAdmin     UserRole = "admin"
	UserRoleModerator UserRole = "moderator"
	UserRoleUser      UserRole = "user"
)

// DefaultRole is assigned to every newly registered user.
const DefaultRole UserRole = UserRoleUser

// User is the core identity entity. It holds only primitive/value types and
// has zero external dependencies (clean domain rule).
type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	DisplayName  string
	AvatarURL    string
	Bio          string
	Status       UserStatus
	Role         UserRole
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
