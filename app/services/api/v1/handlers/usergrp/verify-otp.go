package usergrp

import (
	"context"
	"net/http"
)

type VerifyOTPDTO struct {
	UserID string `json:"phone" validate:"required"`
	OTP    string `json:"otp" validate:"required"`
}

//func toVerifyOTPEntity(d RegisterDTO) (user.RegisterEntity, error) {
//	usr := user.RegisterEntity{
//		FirstName: d.FirstName,
//		LastName:  d.LastName,
//		Phone:     d.Phone,
//		Password:  d.Password,
//	}
//
//	return usr, nil
//}

// VerifyOTP verify user OTP.
func (h *Handlers) VerifyOTP(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// verify user's OTP
	// updated user status to Active
	return nil
}
