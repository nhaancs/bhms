// Package usergrp maintains the group of handlers for user access.
package usergrp

import (
	"errors"
	"github.com/nhaancs/bhms/business/core/user"
	"github.com/nhaancs/bhms/business/web/auth"
	"github.com/nhaancs/bhms/foundation/sms"
	"time"
)

// Set of error variables for handling product group errors.
var (
	ErrInvalidID     = errors.New("id is not in its proper form")
	ErrInvalidStatus = errors.New("status is not allowed for this action")
	ErrInvalidOTP    = errors.New("otp is invalid")
)

// Handlers manages the set of user endpoints.
type Handlers struct {
	user  *user.Core
	auth  *auth.Auth
	sms   *sms.SMS
	keyID string
}

// New constructs a handlers for route access.
func New(
	user *user.Core,
	auth *auth.Auth,
	keyID string,
	sms *sms.SMS,
) *Handlers {
	return &Handlers{
		user:  user,
		auth:  auth,
		keyID: keyID,
		sms:   sms,
	}
}

// UserDTO represents information about an individual user.
type UserDTO struct {
	ID           string   `json:"id"`
	FirstName    string   `json:"first_name"`
	LastName     string   `json:"last_name"`
	Phone        string   `json:"phone"`
	Roles        []string `json:"roles"`
	PasswordHash []byte   `json:"-"`
	Status       string   `json:"status"`
	CreatedAt    string   `json:"CreatedAt"`
	UpdatedAt    string   `json:"UpdatedAt"`
}

func toUserDTO(e user.UserEntity) UserDTO {
	roles := make([]string, len(e.Roles))
	for i, role := range e.Roles {
		roles[i] = role.Name()
	}

	return UserDTO{
		ID:           e.ID.String(),
		FirstName:    e.FirstName,
		LastName:     e.LastName,
		Phone:        e.Phone,
		PasswordHash: e.PasswordHash,
		Roles:        roles,
		Status:       e.Status.Name(),
		CreatedAt:    e.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    e.UpdatedAt.Format(time.RFC3339),
	}
}
