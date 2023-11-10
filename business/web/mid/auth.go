package mid

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/nhaancs/bhms/business/web/response"
	"net/http"

	"github.com/nhaancs/bhms/business/web/auth"
	"github.com/nhaancs/bhms/foundation/web"
)

// Set of error variables for handling user group errors.
var (
	ErrInvalidID = errors.New("ID is not in its proper form")
)

// Authenticate validates a JWT from the `Authorization` header.
func Authenticate(a *auth.Auth) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			claims, err := a.Authenticate(ctx, r.Header.Get("authorization"))
			if err != nil {
				return auth.NewAuthError("authenticate: failed: %s", err)
			}

			if claims.Subject == "" {
				return auth.NewAuthError("authorize: you are not authorized for that action, no claims")
			}

			subjectID, err := uuid.Parse(claims.Subject)
			if err != nil {
				return response.NewError(ErrInvalidID, http.StatusBadRequest)
			}

			ctx = auth.SetUserID(ctx, subjectID)
			ctx = auth.SetClaims(ctx, claims)

			return handler(ctx, w, r)
		}

		return h
	}

	return m
}

// TODO: add authorize middlewares
