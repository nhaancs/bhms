package property

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/nhaancs/bhms/foundation/logger"
	"time"
)

// TODO: limit the number of properties can be created by a user

type Storer interface {
	Create(ctx context.Context, prprty Property) error
	Update(ctx context.Context, prprty Property) error
	Delete(ctx context.Context, prprty Property) error
	QueryByID(ctx context.Context, prprtyID uuid.UUID) (Property, error)
	QueryByManagerID(ctx context.Context, managerID uuid.UUID) ([]Property, error)
}

type Core struct {
	store Storer
	log   *logger.Logger
}

func NewCore(log *logger.Logger, store Storer) *Core {
	return &Core{
		store: store,
		log:   log,
	}
}

func (c *Core) Create(ctx context.Context, e NewProperty) (Property, error) {
	now := time.Now()

	prprty := Property{
		ID:        uuid.New(),
		Status:    e.Status,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := c.store.Create(ctx, prprty); err != nil {
		return Property{}, fmt.Errorf("create: %w", err)
	}

	return prprty, nil
}

//func (c *Core) Update(ctx context.Context, usr User, uu UpdateUser) (User, error) {
//	if uu.FirstName != nil {
//		usr.FirstName = *uu.FirstName
//	}
//
//	if uu.LastName != nil {
//		usr.LastName = *uu.LastName
//	}
//
//	if uu.Phone != nil {
//		usr.Phone = *uu.Phone
//	}
//
//	if uu.Roles != nil {
//		usr.Roles = uu.Roles
//	}
//
//	if uu.Password != nil {
//		pw, err := bcrypt.GenerateFromPassword([]byte(*uu.Password), bcrypt.DefaultCost)
//		if err != nil {
//			return User{}, fmt.Errorf("generatefrompassword: %w", err)
//		}
//		usr.PasswordHash = pw
//	}
//
//	if uu.Status != nil {
//		usr.Status = *uu.Status
//	}
//
//	usr.UpdatedAt = time.Now()
//
//	if err := c.store.Update(ctx, usr); err != nil {
//		return User{}, fmt.Errorf("update: %w", err)
//	}
//
//	return usr, nil
//}

//func (c *Core) QueryByID(ctx context.Context, userID uuid.UUID) (User, error) {
//	user, err := c.store.QueryByID(ctx, userID)
//	if err != nil {
//		return User{}, fmt.Errorf("query: userID[%s]: %w", userID, err)
//	}
//
//	return user, nil
//}

//func (c *Core) QueryByPhone(ctx context.Context, phone string) (User, error) {
//	user, err := c.store.QueryByPhone(ctx, phone)
//	if err != nil {
//		return User{}, fmt.Errorf("query: phone[%s]: %w", phone, err)
//	}
//
//	return user, nil
//}
