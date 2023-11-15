package mid

import (
	"context"
	"github.com/google/uuid"
	"github.com/nhaancs/bhms/business/core/property"
	"github.com/nhaancs/bhms/business/web/auth"
)

// ctxKey represents the type of value for the context key.
type ctxKey int

const (
	claimKey ctxKey = iota + 1
	userIDKey
	propertyKey
)

// =============================================================================

// setClaims stores the claims in the context.
func setClaims(ctx context.Context, claims auth.Claims) context.Context {
	return context.WithValue(ctx, claimKey, claims)
}

// GetClaims returns the claims from the context.
func GetClaims(ctx context.Context) auth.Claims {
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
