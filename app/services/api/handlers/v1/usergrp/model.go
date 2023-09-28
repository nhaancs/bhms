package usergrp

import (
	"fmt"
	"net/mail"
	"time"

	"github.com/nhaancs/realworld/business/core/user"
	"github.com/nhaancs/realworld/foundation/validate"
)

// AppUser represents information about an individual user.
type AppUser struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Bio       string `json:"bio"`
	Image     string `json:"image"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func toAppUser(usr user.User) AppUser {
	return AppUser{
		ID:        usr.ID.String(),
		Username:  usr.Username,
		Email:     usr.Email.Address,
		Bio:       usr.Bio,
		Image:     usr.Image,
		CreatedAt: usr.CreatedAt.Format(time.RFC3339),
		UpdatedAt: usr.UpdatedAt.Format(time.RFC3339),
	}
}

func toAppUsers(users []user.User) []AppUser {
	items := make([]AppUser, len(users))
	for i, usr := range users {
		items[i] = toAppUser(usr)
	}

	return items
}

// =============================================================================

// AppNewUser contains information needed to create a new user.
type AppNewUser struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func toCoreNewUser(app AppNewUser) (user.NewUser, error) {
	addr, err := mail.ParseAddress(app.Email)
	if err != nil {
		return user.NewUser{}, fmt.Errorf("parsing email: %w", err)
	}

	usr := user.NewUser{
		Username: app.Username,
		Email:    *addr,
		Password: app.Password,
	}

	return usr, nil
}

// Validate checks the data in the model is considered clean.
func (app AppNewUser) Validate() error {
	if err := validate.Check(app); err != nil {
		return err
	}

	return nil
}

// =============================================================================

// AppUpdateUser contains information needed to update a user.
type AppUpdateUser struct {
	Email *string `json:"email" validate:"omitempty,email"`
	Bio   *string `json:"bio" validate:"omitempty,email"`
	Image *string `json:"image" validate:"omitempty,email"`
}

func toCoreUpdateUser(app AppUpdateUser) (user.UpdateUser, error) {
	var addr *mail.Address
	if app.Email != nil {
		var err error
		addr, err = mail.ParseAddress(*app.Email)
		if err != nil {
			return user.UpdateUser{}, fmt.Errorf("parsing email: %w", err)
		}
	}

	nu := user.UpdateUser{
		Email: addr,
		Bio:   app.Bio,
		Image: app.Image,
	}

	return nu, nil
}

// Validate checks the data in the model is considered clean.
func (app AppUpdateUser) Validate() error {
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
