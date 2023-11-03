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
		ID:              uuid.New(),
		ManagerID:       e.ManagerID,
		Name:            e.Name,
		AddressLevel1ID: e.AddressLevel1ID,
		AddressLevel2ID: e.AddressLevel2ID,
		AddressLevel3ID: e.AddressLevel3ID,
		Street:          e.Street,
		Status:          e.Status,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := c.store.Create(ctx, prprty); err != nil {
		return Property{}, fmt.Errorf("create: %w", err)
	}

	return prprty, nil
}

func (c *Core) Update(ctx context.Context, prprty Property, up UpdateProperty) (Property, error) {
	if up.Name != nil {
		prprty.Name = *up.Name
	}

	if up.AddressLevel1ID != nil {
		prprty.AddressLevel1ID = *up.AddressLevel1ID
	}

	if up.AddressLevel2ID != nil {
		prprty.AddressLevel2ID = *up.AddressLevel2ID
	}

	if up.AddressLevel3ID != nil {
		prprty.AddressLevel3ID = *up.AddressLevel3ID
	}

	if up.Street != nil {
		prprty.Street = *up.Street
	}

	if up.Status != nil {
		prprty.Status = *up.Status
	}

	prprty.UpdatedAt = time.Now()

	if err := c.store.Update(ctx, prprty); err != nil {
		return Property{}, fmt.Errorf("update: %w", err)
	}

	return prprty, nil
}

func (c *Core) QueryByID(ctx context.Context, propertyID uuid.UUID) (Property, error) {
	prprty, err := c.store.QueryByID(ctx, propertyID)
	if err != nil {
		return Property{}, fmt.Errorf("query: propertyID[%s]: %w", propertyID, err)
	}

	return prprty, nil
}

func (c *Core) QueryByManagerID(ctx context.Context, managerID uuid.UUID) ([]Property, error) {
	prprties, err := c.store.QueryByManagerID(ctx, managerID)
	if err != nil {
		return nil, fmt.Errorf("query: manager id[%s]: %w", managerID.String(), err)
	}

	return prprties, nil
}
