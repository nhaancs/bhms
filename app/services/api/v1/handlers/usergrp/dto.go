package usergrp

import (
	"time"

	"github.com/nhaancs/bhms/business/core/user"
	"github.com/nhaancs/bhms/foundation/validate"
)

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

// =============================================================================

// RegisterDTO contains information needed for a new user to register.
type RegisterDTO struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone" validate:"required,phone"`
	Password  string `json:"password" validate:"required"`
}

func toRegisterEntity(d RegisterDTO) (user.RegisterEntity, error) {
	usr := user.RegisterEntity{
		FirstName: d.FirstName,
		LastName:  d.LastName,
		Phone:     d.Phone,
		Password:  d.Password,
	}

	return usr, nil
}

// Validate checks the data in the model is considered clean.
func (dto RegisterDTO) Validate() error {
	if err := validate.Check(dto); err != nil {
		return err
	}

	return nil
}

// =============================================================================

type token struct {
	Token string `json:"token"`
}

func toToken(v string) token {
	return token{
		Token: v,
	}
}
