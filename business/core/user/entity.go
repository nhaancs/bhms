package user

import (
	"time"

	"github.com/google/uuid"
)

// UserEntity represents information about an individual user.
type UserEntity struct {
	ID           uuid.UUID
	FirstName    string
	LastName     string
	Phone        string
	PasswordHash []byte
	Roles        []Role
	Status       Status
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// RegisterEntity contains information needed to create a new user.
type RegisterEntity struct {
	FirstName string
	LastName  string
	Phone     string
	Password  string
}
