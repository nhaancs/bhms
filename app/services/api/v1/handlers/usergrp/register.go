package usergrp

import (
	"context"
	"errors"
	"fmt"
	"github.com/nhaancs/bhms/app/services/api/v1/request"
	"github.com/nhaancs/bhms/business/core/user"
	"github.com/nhaancs/bhms/foundation/sms"
	"github.com/nhaancs/bhms/foundation/validate"
	"github.com/nhaancs/bhms/foundation/web"
	"net/http"
)

// RegisterDTO contains information needed for a new user to register.
type RegisterDTO struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone" validate:"required,phone"`
	Password  string `json:"password" validate:"required"`
}

func toNewUserEntity(d RegisterDTO) (user.NewUserEntity, error) {
	usr := user.NewUserEntity{
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

// Register adds a new user to the system.
// todo: do rate limit for this api to prevent sending to many sms
func (h *Handlers) Register(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var dto RegisterDTO
	if err := web.Decode(r, &dto); err != nil {
		return request.NewError(err, http.StatusBadRequest)
	}

	e, err := toNewUserEntity(dto)
	if err != nil {
		return request.NewError(err, http.StatusBadRequest)
	}

	usr, err := h.user.Create(ctx, e)
	if err != nil {
		if errors.Is(err, user.ErrUniquePhone) {
			return request.NewError(err, http.StatusConflict)
		}
		return fmt.Errorf("register: usr[%+v]: %+v", usr, err)
	}

	if _, err = h.sms.SendOTP(ctx, sms.OTPInfo{Phone: usr.Phone}); err != nil {
		return fmt.Errorf("senotp: usr[%+v]: %+v", usr, err)
	}

	return web.Respond(ctx, w, toUserDTO(usr), http.StatusCreated)
}
