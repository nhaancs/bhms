package usergrp

import (
	"context"
	"net/http"
)

// VerifyOTP verify user OTP.
func (h *Handlers) VerifyOTP(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// verify user's OTP
	// updated user status to Active
	return nil
}
