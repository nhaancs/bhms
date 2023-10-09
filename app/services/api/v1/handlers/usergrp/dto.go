package usergrp

import (
	"fmt"
	"net/mail"
	"time"

	"github.com/nhaancs/bhms/business/core/user"
	"github.com/nhaancs/bhms/foundation/validate"
)

// UserDTO represents information about an individual user.
type UserDTO struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Email        string   `json:"email"`
	Roles        []string `json:"roles"`
	PasswordHash []byte   `json:"-"`
	Department   string   `json:"department"`
	Enabled      bool     `json:"enabled"`
	DateCreated  string   `json:"dateCreated"`
	DateUpdated  string   `json:"dateUpdated"`
}

func toUserDTO(u user.UserEntity) UserDTO {
	roles := make([]string, len(u.Roles))
	for i, role := range u.Roles {
		roles[i] = role.Name()
	}

	return UserDTO{
		ID:           u.ID.String(),
		Name:         u.Name,
		Email:        u.Email.Address,
		Roles:        roles,
		PasswordHash: u.PasswordHash,
		Department:   u.Department,
		Enabled:      u.Enabled,
		DateCreated:  u.DateCreated.Format(time.RFC3339),
		DateUpdated:  u.DateUpdated.Format(time.RFC3339),
	}
}

func toUserDTOs(users []user.UserEntity) []UserDTO {
	items := make([]UserDTO, len(users))
	for i, usr := range users {
		items[i] = toUserDTO(usr)
	}

	return items
}

// =============================================================================

// RegisterDTO contains information needed to create a new user.
type RegisterDTO struct {
	Name            string   `json:"name" validate:"required"`
	Email           string   `json:"email" validate:"required,email"`
	Roles           []string `json:"roles" validate:"required"`
	Department      string   `json:"department"`
	Password        string   `json:"password" validate:"required"`
	PasswordConfirm string   `json:"passwordConfirm" validate:"eqfield=Password"`
}

func toRegisterEntity(u RegisterDTO) (user.RegisterEntity, error) {
	roles := make([]user.Role, len(u.Roles))
	for i, roleStr := range u.Roles {
		role, err := user.ParseRole(roleStr)
		if err != nil {
			return user.RegisterEntity{}, fmt.Errorf("parsing role: %w", err)
		}
		roles[i] = role
	}

	addr, err := mail.ParseAddress(u.Email)
	if err != nil {
		return user.RegisterEntity{}, fmt.Errorf("parsing email: %w", err)
	}

	usr := user.RegisterEntity{
		Name:            u.Name,
		Email:           *addr,
		Roles:           roles,
		Department:      u.Department,
		Password:        u.Password,
		PasswordConfirm: u.PasswordConfirm,
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

// UpdateUserDTO contains information needed to update a user.
type UpdateUserDTO struct {
	Name            *string  `json:"name"`
	Email           *string  `json:"email" validate:"omitempty,email"`
	Roles           []string `json:"roles"`
	Department      *string  `json:"department"`
	Password        *string  `json:"password"`
	PasswordConfirm *string  `json:"passwordConfirm" validate:"omitempty,eqfield=Password"`
	Enabled         *bool    `json:"enabled"`
}

func toUpdateUserEntity(uu UpdateUserDTO) (user.UpdateUserEntity, error) {
	var roles []user.Role
	if uu.Roles != nil {
		roles = make([]user.Role, len(uu.Roles))
		for i, roleStr := range uu.Roles {
			role, err := user.ParseRole(roleStr)
			if err != nil {
				return user.UpdateUserEntity{}, fmt.Errorf("parsing role: %w", err)
			}
			roles[i] = role
		}
	}

	var addr *mail.Address
	if uu.Email != nil {
		var err error
		addr, err = mail.ParseAddress(*uu.Email)
		if err != nil {
			return user.UpdateUserEntity{}, fmt.Errorf("parsing email: %w", err)
		}
	}

	nu := user.UpdateUserEntity{
		Name:            uu.Name,
		Email:           addr,
		Roles:           roles,
		Department:      uu.Department,
		Password:        uu.Password,
		PasswordConfirm: uu.PasswordConfirm,
		Enabled:         uu.Enabled,
	}

	return nu, nil
}

// Validate checks the data in the model is considered clean.
func (dto UpdateUserDTO) Validate() error {
	if err := validate.Check(dto); err != nil {
		return fmt.Errorf("validate: %w", err)
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
