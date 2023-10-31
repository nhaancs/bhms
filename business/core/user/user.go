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
	ErrInvalidUserStatus     = errors.New("invalid user status")
)

// =============================================================================

// Storer interface declares the behavior this package needs to perists and retrieve data.
type Storer interface {
	Create(ctx context.Context, usr User) error
	Update(ctx context.Context, usr User) error
	Delete(ctx context.Context, usr User) error
	QueryByID(ctx context.Context, userID uuid.UUID) (User, error)
	QueryByIDs(ctx context.Context, userID []uuid.UUID) ([]User, error)
	QueryByPhone(ctx context.Context, phone string) (User, error)
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

// Create a new user to the system.
func (c *Core) Create(ctx context.Context, e NewUser) (User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(e.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("generatefrompassword: %w", err)
	}

	now := time.Now()

	usr := User{
		ID:           uuid.New(),
		FirstName:    e.FirstName,
		LastName:     e.LastName,
		Phone:        e.Phone,
		PasswordHash: hash,
		Roles:        e.Roles,
		Status:       e.Status,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := c.store.Create(ctx, usr); err != nil {
		return User{}, fmt.Errorf("create: %w", err)
	}

	return usr, nil
}

// Update modifies information about a user.
func (c *Core) Update(ctx context.Context, usr User, uu UpdateUser) (User, error) {
	if uu.FirstName != nil {
		usr.FirstName = *uu.FirstName
	}

	if uu.LastName != nil {
		usr.LastName = *uu.LastName
	}

	if uu.Phone != nil {
		usr.Phone = *uu.Phone
	}

	if uu.Roles != nil {
		usr.Roles = uu.Roles
	}

	if uu.Password != nil {
		pw, err := bcrypt.GenerateFromPassword([]byte(*uu.Password), bcrypt.DefaultCost)
		if err != nil {
			return User{}, fmt.Errorf("generatefrompassword: %w", err)
		}
		usr.PasswordHash = pw
	}

	if uu.Status != nil {
		usr.Status = *uu.Status
	}

	usr.UpdatedAt = time.Now()

	if err := c.store.Update(ctx, usr); err != nil {
		return User{}, fmt.Errorf("update: %w", err)
	}

	return usr, nil
}

// QueryByID finds the user by the specified ID.
func (c *Core) QueryByID(ctx context.Context, userID uuid.UUID) (User, error) {
	user, err := c.store.QueryByID(ctx, userID)
	if err != nil {
		return User{}, fmt.Errorf("query: userID[%s]: %w", userID, err)
	}

	return user, nil
}

// QueryByPhone finds the user by a specified user phone.
func (c *Core) QueryByPhone(ctx context.Context, phone string) (User, error) {
	user, err := c.store.QueryByPhone(ctx, phone)
	if err != nil {
		return User{}, fmt.Errorf("query: phone[%s]: %w", phone, err)
	}

	return user, nil
}

// =============================================================================

// Authenticate finds a user by their phone and verifies their password. On
// success it returns a Claims User representing this user. The claims can be
// used to generate a token for future authentication.
func (c *Core) Authenticate(ctx context.Context, phone, password string) (User, error) {
	usr, err := c.QueryByPhone(ctx, phone)
	if err != nil {
		return User{}, fmt.Errorf("query: phone[%s]: %w", phone, err)
	}

	if !usr.Status.Equal(StatusCreated) {
		c.log.Error(ctx, "invalid user status", "got user", usr)
		return User{}, ErrInvalidUserStatus
	}

	err = bcrypt.CompareHashAndPassword(usr.PasswordHash, []byte(password))
	if err != nil {
		return User{}, fmt.Errorf("comparehashandpassword: %w", ErrAuthenticationFailure)
	}

	return usr, nil
}
