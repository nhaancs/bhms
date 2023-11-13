// Package usergrp maintains the group of handlers for user access.
package usergrp

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/nhaancs/bhms/business/core/user"
	"github.com/nhaancs/bhms/business/web/auth"
	"github.com/nhaancs/bhms/business/web/response"
	"github.com/nhaancs/bhms/foundation/sms"
	"github.com/nhaancs/bhms/foundation/web"
	"net/http"
	"time"
)

// Set of error variables for handling product group errors.
var (
	ErrInvalidID     = errors.New("id is not in its proper form")
	ErrInvalidStatus = errors.New("status is not allowed for this action")
	ErrInvalidOTP    = errors.New("otp is invalid")
)

// Handlers manages the set of user endpoints.
type Handlers struct {
	user  *user.Core
	auth  *auth.Auth
	sms   *sms.SMS
	keyID string
}

// New constructs a handlers for route access.
func New(
	user *user.Core,
	auth *auth.Auth,
	keyID string,
	sms *sms.SMS,
) *Handlers {
	return &Handlers{
		user:  user,
		auth:  auth,
		keyID: keyID,
		sms:   sms,
	}
}

// Register adds a new user to the system.
// TODO: limit the number of user can be created (use reCAPTCHA v3, limit by ip, device id)
func (h *Handlers) Register(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app AppRegister
	if err := web.Decode(r, &app); err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	c, err := toCoreNewUser(app)
	if err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	usr, err := h.user.Create(ctx, c)
	if err != nil {
		if errors.Is(err, user.ErrUniquePhone) {
			return response.NewError(err, http.StatusConflict)
		}
		return fmt.Errorf("register: usr[%+v]: %w", app, err)
	}

	//if _, err = h.sms.SendOTP(ctx, sms.OTPInfo{Phone: usr.Phone}); err != nil {
	//	return fmt.Errorf("senotp: usr[%+v]: %w", usr, err)
	//}

	return web.Respond(ctx, w, toAppUser(usr), http.StatusOK)
}

func (h *Handlers) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	userID := uuid.New() // auth.GetUserID(ctx)
	usr, err := h.user.QueryByID(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrNotFound):
			return response.NewError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querybyid: userID[%s]: %w", userID, err)
		}
	}

	var app AppUpdateUser
	if err := web.Decode(r, &app); err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	c := toCoreUpdateUser(app)
	usr, err = h.user.Update(ctx, usr, c)
	if err != nil {
		if errors.Is(err, user.ErrUniquePhone) {
			return response.NewError(err, http.StatusConflict)
		}
		return fmt.Errorf("update: usr[%+v]: %w", app, err)
	}

	return web.Respond(ctx, w, toAppUser(usr), http.StatusOK)
}

// VerifyOTP verify user OTP.
func (h *Handlers) VerifyOTP(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app AppVerifyOTP
	if err := web.Decode(r, &app); err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	userID, err := uuid.Parse(app.UserID)
	if err != nil {
		return response.NewError(ErrInvalidID, http.StatusBadRequest)
	}

	usr, err := h.user.QueryByID(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrNotFound):
			return response.NewError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querybyid: userID[%s]: %w", userID, err)
		}
	}
	if !usr.Status.Equal(user.StatusCreated) {
		return response.NewError(ErrInvalidStatus, http.StatusBadRequest)
	}

	err = h.sms.CheckOTP(ctx, sms.VerifyOTPInfo{
		Phone: usr.Phone,
		Code:  app.OTP,
	})
	if err != nil {
		return response.NewError(ErrInvalidOTP, http.StatusBadRequest)
	}

	status := user.StatusCreated
	usr, err = h.user.Update(ctx, usr, user.UpdateUser{Status: &status})
	if err != nil {
		return fmt.Errorf("update: userID[%s] app[%+v]: %w", userID, app, err)
	}

	return web.Respond(ctx, w, toAppUser(usr), http.StatusOK)
}

// Token provides an API token for the authenticated user.
func (h *Handlers) Token(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	phone, pass, ok := r.BasicAuth()
	if !ok {
		return auth.NewAuthError("must provide email and password in Basic auth")
	}

	usr, err := h.user.Authenticate(ctx, phone, pass)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrNotFound):
			return response.NewError(err, http.StatusNotFound)
		case errors.Is(err, user.ErrAuthenticationFailure) || errors.Is(err, user.ErrInvalidUserStatus):
			return auth.NewAuthError(err.Error())
		default:
			return fmt.Errorf("authenticate: %w", err)
		}
	}

	claims := auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   usr.ID.String(),
			Issuer:    "bhms",
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(30 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Roles: usr.Roles,
	}

	token, err := h.auth.GenerateToken(ctx, h.keyID, claims)
	if err != nil {
		return fmt.Errorf("generatetoken: %w", err)
	}

	return web.Respond(ctx, w, toToken(token), http.StatusOK)
}
