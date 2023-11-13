package mid

import (
	"context"
	"github.com/google/uuid"
	"github.com/nhaancs/bhms/business/web/auth"
)

// ctxKey represents the type of value for the context key.
type ctxKey int

// key is used to store/retrieve a Claims value from a context.Context.
const claimKey ctxKey = 1

// key is used to store/retrieve a user value from a context.Context.
const userKey ctxKey = 2

// =============================================================================

// setClaims stores the claims in the context.
func setClaims(ctx context.Context, claims auth.Claims) context.Context {
	return context.WithValue(ctx, claimKey, claims)
}

// getClaims returns the claims from the context.
func getClaims(ctx context.Context) auth.Claims {
	v, ok := ctx.Value(claimKey).(auth.Claims)
	if !ok {
		return auth.Claims{}
	}
	return v
}

// =============================================================================

// setUserID stores the user id from the request in the context.
func setUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userKey, userID)
}

// getUserID returns the claims from the context.
func getUserID(ctx context.Context) uuid.UUID {
	v, ok := ctx.Value(userKey).(uuid.UUID)
	if !ok {
		return uuid.UUID{}
	}
	return v
}
