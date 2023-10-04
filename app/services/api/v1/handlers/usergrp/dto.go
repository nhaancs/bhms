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

func toUserDTO(usr user.UserEntity) UserDTO {
	roles := make([]string, len(usr.Roles))
	for i, role := range usr.Roles {
		roles[i] = role.Name()
	}

	return UserDTO{
		ID:           usr.ID.String(),
		Name:         usr.Name,
		Email:        usr.Email.Address,
		Roles:        roles,
		PasswordHash: usr.PasswordHash,
		Department:   usr.Department,
		Enabled:      usr.Enabled,
		DateCreated:  usr.DateCreated.Format(time.RFC3339),
		DateUpdated:  usr.DateUpdated.Format(time.RFC3339),
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

// NewUserDTO contains information needed to create a new user.
type NewUserDTO struct {
	Name            string   `json:"name" validate:"required"`
	Email           string   `json:"email" validate:"required,email"`
	Roles           []string `json:"roles" validate:"required"`
	Department      string   `json:"department"`
	Password        string   `json:"password" validate:"required"`
	PasswordConfirm string   `json:"passwordConfirm" validate:"eqfield=Password"`
}

func toNewUserEntity(u NewUserDTO) (user.NewUserEntity, error) {
	roles := make([]user.Role, len(u.Roles))
	for i, roleStr := range u.Roles {
		role, err := user.ParseRole(roleStr)
		if err != nil {
			return user.NewUserEntity{}, fmt.Errorf("parsing role: %w", err)
		}
		roles[i] = role
	}

	addr, err := mail.ParseAddress(u.Email)
	if err != nil {
		return user.NewUserEntity{}, fmt.Errorf("parsing email: %w", err)
	}

	usr := user.NewUserEntity{
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
func (app NewUserDTO) Validate() error {
	if err := validate.Check(app); err != nil {
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

func toUpdateUserEntity(app UpdateUserDTO) (user.UpdateUserEntity, error) {
	var roles []user.Role
	if app.Roles != nil {
		roles = make([]user.Role, len(app.Roles))
		for i, roleStr := range app.Roles {
			role, err := user.ParseRole(roleStr)
			if err != nil {
				return user.UpdateUserEntity{}, fmt.Errorf("parsing role: %w", err)
			}
			roles[i] = role
		}
	}

	var addr *mail.Address
	if app.Email != nil {
		var err error
		addr, err = mail.ParseAddress(*app.Email)
		if err != nil {
			return user.UpdateUserEntity{}, fmt.Errorf("parsing email: %w", err)
		}
	}

	nu := user.UpdateUserEntity{
		Name:            app.Name,
		Email:           addr,
		Roles:           roles,
		Department:      app.Department,
		Password:        app.Password,
		PasswordConfirm: app.PasswordConfirm,
		Enabled:         app.Enabled,
	}

	return nu, nil
}

// Validate checks the data in the model is considered clean.
func (app UpdateUserDTO) Validate() error {
	if err := validate.Check(app); err != nil {
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
