package property

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/nhaancs/bhms/foundation/logger"
	"time"
)

var (
	ErrNotFound      = errors.New("property not found")
	ErrLimitExceeded = errors.New("max number of properties exceeded")
)

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

// TODO: limit the number of blocks, floors, units can be created
func (c *Core) Create(ctx context.Context, core NewProperty) (Property, error) {
	prprties, err := c.store.QueryByManagerID(ctx, core.ManagerID)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return Property{}, err
	}
	if len(prprties) > 0 {
		return Property{}, ErrLimitExceeded
	}

	now := time.Now()
	prprty := Property{
		ID:              uuid.New(),
		ManagerID:       core.ManagerID,
		Name:            core.Name,
		AddressLevel1ID: core.AddressLevel1ID,
		AddressLevel2ID: core.AddressLevel2ID,
		AddressLevel3ID: core.AddressLevel3ID,
		Street:          core.Street,
		Status:          core.Status,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := c.store.Create(ctx, prprty); err != nil {
		return Property{}, fmt.Errorf("create: %w", err)
	}

	return prprty, nil
}

func (c *Core) Update(ctx context.Context, o Property, n UpdateProperty) (Property, error) {
	if n.Name != nil {
		o.Name = *n.Name
	}

	if n.AddressLevel1ID != nil {
		o.AddressLevel1ID = *n.AddressLevel1ID
	}

	if n.AddressLevel2ID != nil {
		o.AddressLevel2ID = *n.AddressLevel2ID
	}

	if n.AddressLevel3ID != nil {
		o.AddressLevel3ID = *n.AddressLevel3ID
	}

	if n.Street != nil {
		o.Street = *n.Street
	}

	if n.Status != nil {
		o.Status = *n.Status
	}

	o.UpdatedAt = time.Now()

	if err := c.store.Update(ctx, o); err != nil {
		return Property{}, fmt.Errorf("update: %w", err)
	}

	return o, nil
}

func (c *Core) QueryByID(ctx context.Context, id uuid.UUID) (Property, error) {
	prprty, err := c.store.QueryByID(ctx, id)
	if err != nil {
		return Property{}, fmt.Errorf("query: id[%s]: %w", id, err)
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
