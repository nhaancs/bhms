package user

import (
	"time"

	"github.com/google/uuid"
)

// User represents information about an individual user.
type User struct {
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

// NewUser contains information needed to create a new user.
type NewUser struct {
	FirstName string
	LastName  string
	Phone     string
	Password  string
}

// UpdateUser contains information needed to update a user.
type UpdateUser struct {
	FirstName *string
	LastName  *string
	Phone     *string
	Roles     []Role
	Password  *string
	Status    *Status
}
