package usergrp

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/nhaancs/bhms/app/services/api/v1/request"
	"github.com/nhaancs/bhms/business/core/user"
	"github.com/nhaancs/bhms/foundation/sms"
	"github.com/nhaancs/bhms/foundation/validate"
	"github.com/nhaancs/bhms/foundation/web"
	"net/http"
)

type AppVerifyOTP struct {
	UserID string `json:"user_id" validate:"required"`
	OTP    string `json:"otp" validate:"required"`
}

// Validate checks the data in the model is considered clean.
func (r AppVerifyOTP) Validate() error {
	if err := validate.Check(r); err != nil {
		return err
	}

	return nil
}

// VerifyOTP verify user OTP.
func (h *Handlers) VerifyOTP(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app AppVerifyOTP
	if err := web.Decode(r, &app); err != nil {
		return request.NewError(err, http.StatusBadRequest)
	}

	userID, err := uuid.Parse(app.UserID)
	if err != nil {
		return request.NewError(ErrInvalidID, http.StatusBadRequest)
	}

	usr, err := h.user.QueryByID(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrNotFound):
			return request.NewError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querybyid: userID[%s]: %+v", userID, err)
		}
	}
	if usr.Status != user.StatusCreated {
		return request.NewError(ErrInvalidStatus, http.StatusBadRequest)
	}

	err = h.sms.CheckOTP(ctx, sms.VerifyOTPInfo{
		Phone: usr.Phone,
		Code:  app.OTP,
	})
	if err != nil {
		return request.NewError(ErrInvalidOTP, http.StatusBadRequest)
	}

	status := user.StatusCreated
	usr, err = h.user.Update(ctx, usr, user.UpdateUser{Status: &status})
	if err != nil {
		return fmt.Errorf("update: userID[%s] app[%+v]: %w", userID, app, err)
	}

	return web.Respond(ctx, w, toAppUser(usr), http.StatusOK)
}
