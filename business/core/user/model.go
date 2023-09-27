package user

import (
	"net/mail"
	"time"

	"github.com/google/uuid"
)

// User represents information about an individual user.
type User struct {
	ID           uuid.UUID
	Username     string
	Email        mail.Address
	Bio          string
	Image        string
	PasswordHash []byte
	DateCreated  time.Time
	DateUpdated  time.Time
}

// NewUser contains information needed to create a new user.
type NewUser struct {
	Username string
	Email    mail.Address
	Password string
}

// UpdateUser contains information needed to update a user.
type UpdateUser struct {
	Email *mail.Address
	Bio   *string
	Image *string
}
