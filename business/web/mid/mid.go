package mid

import (
	"context"
	"github.com/google/uuid"
	"github.com/nhaancs/bhms/business/core/property"
	"github.com/nhaancs/bhms/business/core/user"
	"github.com/nhaancs/bhms/business/web/auth"
)

// ctxKey represents the type of value for the context key.
type ctxKey int

const (
	claimKey ctxKey = iota + 1
	userIDKey
	userKey
	propertyKey
)

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

func setUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func GetUserID(ctx context.Context) uuid.UUID {
	v, ok := ctx.Value(userIDKey).(uuid.UUID)
	if !ok {
		return uuid.UUID{}
	}
	return v
}

// =============================================================================

func setProperty(ctx context.Context, prprty property.Property) context.Context {
	return context.WithValue(ctx, propertyKey, prprty)
}

func GetProperty(ctx context.Context) property.Property {
	v, ok := ctx.Value(propertyKey).(property.Property)
	if !ok {
		return property.Property{}
	}
	return v
}

// =============================================================================

func setUser(ctx context.Context, usr user.User) context.Context {
	return context.WithValue(ctx, userKey, usr)
}

func GetUser(ctx context.Context) user.User {
	v, ok := ctx.Value(userKey).(user.User)
	if !ok {
		return user.User{}
	}
	return v
}
