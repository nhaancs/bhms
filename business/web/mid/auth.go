package mid

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/nhaancs/bhms/business/core/property"
	"github.com/nhaancs/bhms/business/core/user"
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

			ctx = setUserID(ctx, subjectID)
			ctx = setClaims(ctx, claims)

			return handler(ctx, w, r)
		}

		return h
	}

	return m
}

// TODO: add authorize middlewares

// Authorize executes the specified role and does not extract any domain data.
func Authorize(a *auth.Auth, rule string) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			claims := getClaims(ctx)
			if err := a.Authorize(ctx, claims, uuid.UUID{}, rule); err != nil {
				return auth.NewAuthError("authorize: you are not authorized for that action, claims[%v] rule[%v]: %s", claims.Roles, rule, err)
			}

			return handler(ctx, w, r)
		}

		return h
	}

	return m
}

// AuthorizeUser executes the specified role and extracts the specified user from the DB if a user id is specified in the call.
// Depending on the rule specified, the userid from the claims may be compared with the specified user id.
func AuthorizeUser(a *auth.Auth, rule string, usrCore *user.Core) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			var userID uuid.UUID

			if id := web.Param(r, "id"); id != "" {
				var err error
				userID, err = uuid.Parse(id)
				if err != nil {
					return response.NewError(ErrInvalidID, http.StatusBadRequest)
				}

				usr, err := usrCore.QueryByID(ctx, userID)
				if err != nil {
					switch {
					case errors.Is(err, user.ErrNotFound):
						return response.NewError(err, http.StatusNoContent)
					default:
						return fmt.Errorf("querybyid: userID[%s]: %w", userID, err)
					}
				}

				ctx = setUser(ctx, usr)
			}

			claims := getClaims(ctx)
			if err := a.Authorize(ctx, claims, userID, rule); err != nil {
				return auth.NewAuthError("authorize: you are not authorized for that action, claims[%v] rule[%v]: %s", claims.Roles, rule, err)
			}

			return handler(ctx, w, r)
		}

		return h
	}

	return m
}

// AuthorizeProperty executes the specified role and extracts the specified property from the DB if a property id is specified in the call.
// Depending on the rule specified, the userid from the claims may be compared with the specified manager id from the property.
func AuthorizeProperty(a *auth.Auth, rule string, prprtyCore *property.Core) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			var userID uuid.UUID

			if id := web.Param(r, "id"); id != "" {
				var err error
				propertyID, err := uuid.Parse(id)
				if err != nil {
					return response.NewError(ErrInvalidID, http.StatusBadRequest)
				}

				prprty, err := prprtyCore.QueryByID(ctx, propertyID)
				if err != nil {
					switch {
					case errors.Is(err, property.ErrNotFound):
						return response.NewError(err, http.StatusNoContent)
					default:
						return fmt.Errorf("querybyid: id[%s]: %w", propertyID, err)
					}
				}

				userID = prprty.ManagerID
				ctx = setProperty(ctx, prprty)
			}

			claims := getClaims(ctx)
			if err := a.Authorize(ctx, claims, userID, rule); err != nil {
				return auth.NewAuthError("authorize: you are not authorized for that action, claims[%v] rule[%v]: %s", claims.Roles, rule, err)
			}

			return handler(ctx, w, r)
		}

		return h
	}

	return m
}
