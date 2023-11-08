package unit

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/nhaancs/bhms/foundation/logger"
	"time"
)

var (
	ErrNotFound = errors.New("unit not found")
)

type Storer interface {
	Create(ctx context.Context, core Unit) error
	BatchCreate(ctx context.Context, cores []Unit) error
	Update(ctx context.Context, core Unit) error
	Delete(ctx context.Context, core Unit) error
	QueryByID(ctx context.Context, id uuid.UUID) (Unit, error)
	QueryByFloorID(ctx context.Context, id uuid.UUID) ([]Unit, error)
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

func (c *Core) Create(ctx context.Context, core NewUnit) (Unit, error) {
	now := time.Now()
	unt := Unit{
		ID:         uuid.New(),
		Name:       core.Name,
		PropertyID: core.PropertyID,
		BlockID:    core.BlockID,
		FloorID:    core.FloorID,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := c.store.Create(ctx, unt); err != nil {
		return Unit{}, fmt.Errorf("create: %w", err)
	}

	return unt, nil
}

func (c *Core) BatchCreate(ctx context.Context, cores []NewUnit) ([]Unit, error) {
	now := time.Now()
	unts := make([]Unit, len(cores))

	for i := range cores {
		unts[i] = Unit{
			ID:         uuid.New(),
			Name:       cores[i].Name,
			PropertyID: cores[i].PropertyID,
			BlockID:    cores[i].BlockID,
			FloorID:    cores[i].FloorID,
			CreatedAt:  now,
			UpdatedAt:  now,
		}
	}

	if err := c.store.BatchCreate(ctx, unts); err != nil {
		return nil, fmt.Errorf("batch create: %w", err)
	}

	return unts, nil
}

func (c *Core) Update(ctx context.Context, o Unit, n UpdateUnit) (Unit, error) {
	if n.Name != nil {
		o.Name = *n.Name
	}
	o.UpdatedAt = time.Now()

	if err := c.store.Update(ctx, o); err != nil {
		return Unit{}, fmt.Errorf("update: %w", err)
	}

	return o, nil
}

func (c *Core) QueryByID(ctx context.Context, id uuid.UUID) (Unit, error) {
	unt, err := c.store.QueryByID(ctx, id)
	if err != nil {
		return Unit{}, fmt.Errorf("query: id[%s]: %w", id, err)
	}

	return unt, nil
}

func (c *Core) QueryByFloorID(ctx context.Context, id uuid.UUID) ([]Unit, error) {
	unts, err := c.store.QueryByFloorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("query: block id[%s]: %w", id.String(), err)
	}

	return unts, nil
}
