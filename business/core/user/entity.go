package user

import (
	"net/mail"
	"time"

	"github.com/google/uuid"
)

// UserEntity represents information about an individual user.
type UserEntity struct {
	ID           uuid.UUID
	Name         string
	Email        mail.Address
	Roles        []Role
	PasswordHash []byte
	Department   string
	Enabled      bool
	DateCreated  time.Time
	DateUpdated  time.Time
}

// RegisterEntity contains information needed to create a new user.
type RegisterEntity struct {
	Name            string
	Email           mail.Address
	Roles           []Role
	Department      string
	Password        string
	PasswordConfirm string
}
