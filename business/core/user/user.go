// Package user provides an example of a core business API. Right now these
// calls are just wrapping the data/data layer. But at some point you will
// want auditing or something that isn't specific to the data/store layer.
package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nhaancs/bhms/foundation/logger"
	"golang.org/x/crypto/bcrypt"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound              = errors.New("user not found")
	ErrUniquePhone           = errors.New("phone is not unique")
	ErrAuthenticationFailure = errors.New("authentication failed")
)

// =============================================================================

// Storer interface declares the behavior this package needs to perists and retrieve data.
type Storer interface {
	Create(ctx context.Context, usr UserEntity) error
	Update(ctx context.Context, usr UserEntity) error
	Delete(ctx context.Context, usr UserEntity) error
	QueryByID(ctx context.Context, userID uuid.UUID) (UserEntity, error)
	QueryByIDs(ctx context.Context, userID []uuid.UUID) ([]UserEntity, error)
	QueryByPhone(ctx context.Context, phone string) (UserEntity, error)
}

// =============================================================================

// Core manages the set of APIs for user access.
type Core struct {
	store Storer
	log   *logger.Logger
}

// NewCore constructs a core for user api access.
func NewCore(log *logger.Logger, store Storer) *Core {
	return &Core{
		store: store,
		log:   log,
	}
}

// Register a new user to the system.
func (c *Core) Register(ctx context.Context, e RegisterEntity) (UserEntity, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(e.Password), bcrypt.DefaultCost)
	if err != nil {
		return UserEntity{}, fmt.Errorf("generatefrompassword: %w", err)
	}

	now := time.Now()

	usr := UserEntity{
		ID:           uuid.New(),
		FirstName:    e.FirstName,
		LastName:     e.LastName,
		Phone:        e.Phone,
		PasswordHash: hash,
		Roles:        []Role{RoleUser},
		Status:       StatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := c.store.Create(ctx, usr); err != nil {
		return UserEntity{}, fmt.Errorf("create: %w", err)
	}

	return usr, nil
}

// QueryByID finds the user by the specified ID.
func (c *Core) QueryByID(ctx context.Context, userID uuid.UUID) (UserEntity, error) {
	user, err := c.store.QueryByID(ctx, userID)
	if err != nil {
		return UserEntity{}, fmt.Errorf("query: userID[%s]: %w", userID, err)
	}

	return user, nil
}

// QueryByPhone finds the user by a specified user phone.
func (c *Core) QueryByPhone(ctx context.Context, phone string) (UserEntity, error) {
	user, err := c.store.QueryByPhone(ctx, phone)
	if err != nil {
		return UserEntity{}, fmt.Errorf("query: phone[%s]: %w", phone, err)
	}

	return user, nil
}

// =============================================================================

// Authenticate finds a user by their phone and verifies their password. On
// success it returns a Claims UserEntity representing this user. The claims can be
// used to generate a token for future authentication.
func (c *Core) Authenticate(ctx context.Context, phone, password string) (UserEntity, error) {
	usr, err := c.QueryByPhone(ctx, phone)
	if err != nil {
		return UserEntity{}, fmt.Errorf("query: phone[%s]: %w", phone, err)
	}

	if err := bcrypt.CompareHashAndPassword(usr.PasswordHash, []byte(password)); err != nil {
		return UserEntity{}, fmt.Errorf("comparehashandpassword: %w", ErrAuthenticationFailure)
	}

	return usr, nil
}
