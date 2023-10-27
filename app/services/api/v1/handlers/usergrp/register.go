package usergrp

import (
	"context"
	"errors"
	"fmt"
	"github.com/nhaancs/bhms/app/services/api/v1/request"
	"github.com/nhaancs/bhms/business/core/user"
	"github.com/nhaancs/bhms/foundation/validate"
	"github.com/nhaancs/bhms/foundation/web"
	"net/http"
)

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

// Register adds a new user to the system.
// TODO:
// - Verify phone number by sending otp
// - Rate limit for this api to prevent sending to many sms
func (h *Handlers) Register(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app AppRegister
	if err := web.Decode(r, &app); err != nil {
		return request.NewError(err, http.StatusBadRequest)
	}

	e, err := toCoreNewUser(app)
	if err != nil {
		return request.NewError(err, http.StatusBadRequest)
	}

	usr, err := h.user.Create(ctx, e)
	if err != nil {
		if errors.Is(err, user.ErrUniquePhone) {
			return request.NewError(err, http.StatusConflict)
		}
		return fmt.Errorf("register: usr[%+v]: %w", app, err)
	}

	//if _, err = h.sms.SendOTP(ctx, sms.OTPInfo{Phone: usr.Phone}); err != nil {
	//	return fmt.Errorf("senotp: usr[%+v]: %w", usr, err)
	//}

	return web.Respond(ctx, w, toAppUser(usr), http.StatusCreated)
}
