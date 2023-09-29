// Package usergrp maintains the group of handlers for user access.
package usergrp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/nhaancs/realworld/business/core/user"
	"github.com/nhaancs/realworld/business/data/transaction"
	"github.com/nhaancs/realworld/business/web/auth"
	v1 "github.com/nhaancs/realworld/business/web/v1"
	"github.com/nhaancs/realworld/foundation/validate"
	"github.com/nhaancs/realworld/foundation/web"
)

// Handlers manages the set of user endpoints.
type Handlers struct {
	user *user.Core
	auth *auth.Auth
}

// New constructs a handlers for route access.
func New(user *user.Core, auth *auth.Auth) *Handlers {
	return &Handlers{
		user: user,
		auth: auth,
	}
}

// executeUnderTransaction constructs a new Handlers value with the core apis
// using a store transaction that was created via middleware.
func (h *Handlers) executeUnderTransaction(ctx context.Context) (*Handlers, error) {
	if tx, ok := transaction.Get(ctx); ok {
		user, err := h.user.ExecuteUnderTransaction(tx)
		if err != nil {
			return nil, err
		}

		h = &Handlers{
			user: user,
			auth: h.auth,
		}

		return h, nil
	}

	return h, nil
}

// Create adds a new user to the system.
func (h *Handlers) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app AppNewUser
	if err := web.Decode(r, &app); err != nil {
		return v1.NewRequestError(err, http.StatusBadRequest)
	}

	nc, err := toCoreNewUser(app)
	if err != nil {
		return v1.NewRequestError(err, http.StatusBadRequest)
	}

	usr, err := h.user.Create(ctx, nc)
	if err != nil {
		if errors.Is(err, user.ErrUniqueEmail) {
			return v1.NewRequestError(err, http.StatusConflict)
		}
		return fmt.Errorf("create: usr[%+v]: %w", usr, err)
	}

	return web.Respond(ctx, w, toAppUser(usr), http.StatusCreated)
}

// Update updates a user in the system.
func (h *Handlers) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	h, err := h.executeUnderTransaction(ctx)
	if err != nil {
		return err
	}

	var app AppUpdateUser
	if err := web.Decode(r, &app); err != nil {
		return v1.NewRequestError(err, http.StatusBadRequest)
	}

	userID := auth.GetUserID(ctx)

	usr, err := h.user.QueryByID(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrNotFound):
			return v1.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querybyid: userID[%s]: %w", userID, err)
		}
	}

	uu, err := toCoreUpdateUser(app)
	if err != nil {
		return v1.NewRequestError(err, http.StatusBadRequest)
	}

	usr, err = h.user.Update(ctx, usr, uu)
	if err != nil {
		return fmt.Errorf("update: userID[%s] uu[%+v]: %w", userID, uu, err)
	}

	return web.Respond(ctx, w, toAppUser(usr), http.StatusOK)
}

// QueryByID returns a user by its ID.
func (h *Handlers) QueryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := auth.GetUserID(ctx)

	usr, err := h.user.QueryByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrNotFound):
			return v1.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querybyid: id[%s]: %w", id, err)
		}
	}

	return web.Respond(ctx, w, toAppUser(usr), http.StatusOK)
}

// Token provides an API token for the authenticated user.
func (h *Handlers) Token(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	kid := web.Param(r, "kid")
	if kid == "" {
		return validate.NewFieldsError("kid", errors.New("missing kid"))
	}

	email, pass, ok := r.BasicAuth()
	if !ok {
		return auth.NewAuthError("must provide email and password in Basic auth")
	}

	addr, err := mail.ParseAddress(email)
	if err != nil {
		return auth.NewAuthError("invalid email format")
	}

	usr, err := h.user.Authenticate(ctx, *addr, pass)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrNotFound):
			return v1.NewRequestError(err, http.StatusNotFound)
		case errors.Is(err, user.ErrAuthenticationFailure):
			return auth.NewAuthError(err.Error())
		default:
			return fmt.Errorf("authenticate: %w", err)
		}
	}

	claims := auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   usr.ID.String(),
			Issuer:    "service project",
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
	}

	token, err := h.auth.GenerateToken(kid, claims)
	if err != nil {
		return fmt.Errorf("generatetoken: %w", err)
	}

	return web.Respond(ctx, w, toToken(token), http.StatusOK)
}
