package usergrp

import (
	"github.com/nhaancs/bhms/business/core/user"
	"github.com/nhaancs/bhms/foundation/validate"
	"time"
)

type AppVerifyOTP struct {
	UserID string `json:"user_id" validate:"required,uuid"`
	OTP    string `json:"otp" validate:"required,len=6"`
}

// Validate checks the data in the model is considered clean.
func (r AppVerifyOTP) Validate() error {
	if err := validate.Check(r); err != nil {
		return err
	}

	return nil
}

// ===============================================================

// AppRegister contains information needed for a new user to register.
type AppRegister struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone" validate:"required,number,startswith=0,len=10"`
	Password  string `json:"password" validate:"required,min=6"`
}

func toCoreNewUser(a AppRegister) (user.NewUser, error) {
	usr := user.NewUser{
		FirstName: a.FirstName,
		LastName:  a.LastName,
		Phone:     a.Phone,
		Password:  a.Password,
		Status:    user.StatusCreated,
		Roles:     []user.Role{user.RoleUser},
	}

	return usr, nil
}

// Validate checks the data in the model is considered clean.
func (r AppRegister) Validate() error {
	if err := validate.Check(r); err != nil {
		return err
	}

	return nil
}

// ===============================================================

type AppUpdateUser struct {
	FirstName *string `json:"first_name" validate:"required"`
	LastName  *string `json:"last_name"`
	Phone     *string `json:"phone" validate:"required,number,startswith=0,len=10"`
	Password  *string `json:"password" validate:"required,min=6"`
}

func toCoreUpdateUser(a AppUpdateUser) (user.UpdateUser, error) {
	usr := user.UpdateUser{
		FirstName: a.FirstName,
		LastName:  a.LastName,
		Phone:     a.Phone,
		Password:  a.Password,
	}

	return usr, nil
}

// Validate checks the data in the model is considered clean.
func (r AppUpdateUser) Validate() error {
	if err := validate.Check(r); err != nil {
		return err
	}

	return nil
}

// ===============================================================

type token struct {
	Token string `json:"token"`
}

func toToken(v string) token {
	return token{
		Token: v,
	}
}

// ===============================================================

// AppUser represents information about an individual user.
type AppUser struct {
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

func toAppUser(e user.User) AppUser {
	roles := make([]string, len(e.Roles))
	for i, role := range e.Roles {
		roles[i] = role.Name()
	}

	return AppUser{
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
