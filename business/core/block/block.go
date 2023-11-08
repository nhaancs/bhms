package block

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/nhaancs/bhms/foundation/logger"
	"time"
)

var (
	ErrNotFound = errors.New("block not found")
)

type Storer interface {
	Create(ctx context.Context, core Block) error
	BatchCreate(ctx context.Context, cores []Block) error
	Update(ctx context.Context, core Block) error
	Delete(ctx context.Context, core Block) error
	QueryByID(ctx context.Context, id uuid.UUID) (Block, error)
	QueryByPropertyID(ctx context.Context, id uuid.UUID) ([]Block, error)
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

func (c *Core) Create(ctx context.Context, core NewBlock) (Block, error) {
	now := time.Now()
	blck := Block{
		ID:         uuid.New(),
		Name:       core.Name,
		PropertyID: core.PropertyID,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := c.store.Create(ctx, blck); err != nil {
		return Block{}, fmt.Errorf("create: %w", err)
	}

	return blck, nil
}

func (c *Core) BatchCreate(ctx context.Context, cores []NewBlock) ([]Block, error) {
	now := time.Now()
	blcks := make([]Block, len(cores))

	for i := range cores {
		blcks[i] = Block{
			ID:         uuid.New(),
			Name:       cores[i].Name,
			PropertyID: cores[i].PropertyID,
			CreatedAt:  now,
			UpdatedAt:  now,
		}
	}

	if err := c.store.BatchCreate(ctx, blcks); err != nil {
		return nil, fmt.Errorf("batch create: %w", err)
	}

	return blcks, nil
}

func (c *Core) Update(ctx context.Context, o Block, n UpdateBlock) (Block, error) {
	if n.Name != nil {
		o.Name = *n.Name
	}
	o.UpdatedAt = time.Now()

	if err := c.store.Update(ctx, o); err != nil {
		return Block{}, fmt.Errorf("update: %w", err)
	}

	return o, nil
}

func (c *Core) QueryByID(ctx context.Context, id uuid.UUID) (Block, error) {
	blck, err := c.store.QueryByID(ctx, id)
	if err != nil {
		return Block{}, fmt.Errorf("query: id[%s]: %w", id, err)
	}

	return blck, nil
}

func (c *Core) QueryByPropertyID(ctx context.Context, id uuid.UUID) ([]Block, error) {
	blcks, err := c.store.QueryByPropertyID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("query: property id[%s]: %w", id.String(), err)
	}

	return blcks, nil
}
