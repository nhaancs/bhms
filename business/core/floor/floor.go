package floor

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/nhaancs/bhms/business/data/transaction"
	"github.com/nhaancs/bhms/foundation/logger"
	"time"
)

var (
	ErrNotFound = errors.New("floor not found")
)

type Storer interface {
	Create(ctx context.Context, core Floor) error
	BatchCreate(ctx context.Context, cores []Floor) error
	Update(ctx context.Context, core Floor) error
	QueryByID(ctx context.Context, id uuid.UUID) (Floor, error)
	QueryByBlockID(ctx context.Context, id uuid.UUID) ([]Floor, error)
	QueryByPropertyID(ctx context.Context, id uuid.UUID) ([]Floor, error)
	ExecuteUnderTransaction(tx transaction.Transaction) (Storer, error)
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

func (c *Core) Create(ctx context.Context, core NewFloor) (Floor, error) {
	now := time.Now()
	flr := Floor{
		ID:         core.ID,
		Name:       core.Name,
		PropertyID: core.PropertyID,
		BlockID:    core.BlockID,
		Status:     StatusActive,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := c.store.Create(ctx, flr); err != nil {
		return Floor{}, fmt.Errorf("create: %w", err)
	}

	return flr, nil
}

func (c *Core) BatchCreate(ctx context.Context, cores []NewFloor) ([]Floor, error) {
	now := time.Now()
	flrs := make([]Floor, len(cores))

	for i := range cores {
		flrs[i] = Floor{
			ID:         uuid.New(),
			Name:       cores[i].Name,
			PropertyID: cores[i].PropertyID,
			BlockID:    cores[i].BlockID,
			Status:     StatusActive,
			CreatedAt:  now,
			UpdatedAt:  now,
		}
	}

	if err := c.store.BatchCreate(ctx, flrs); err != nil {
		return nil, fmt.Errorf("batch create: %w", err)
	}

	return flrs, nil
}

func (c *Core) Update(ctx context.Context, o Floor, n UpdateFloor) (Floor, error) {
	if n.Name != nil {
		o.Name = *n.Name
	}

	o.UpdatedAt = time.Now()

	if err := c.store.Update(ctx, o); err != nil {
		return Floor{}, fmt.Errorf("update: %w", err)
	}

	return o, nil
}

func (c *Core) Delete(ctx context.Context, core Floor) (Floor, error) {
	core.UpdatedAt = time.Now()
	core.Status = StatusDeleted

	if err := c.store.Update(ctx, core); err != nil {
		return Floor{}, fmt.Errorf("update: %w", err)
	}

	return core, nil
}

func (c *Core) QueryByID(ctx context.Context, id uuid.UUID) (Floor, error) {
	flr, err := c.store.QueryByID(ctx, id)
	if err != nil {
		return Floor{}, fmt.Errorf("query: id[%s]: %w", id, err)
	}

	return flr, nil
}

func (c *Core) QueryByBlockID(ctx context.Context, id uuid.UUID) ([]Floor, error) {
	flrs, err := c.store.QueryByBlockID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("query: block id[%s]: %w", id.String(), err)
	}

	return flrs, nil
}

func (c *Core) QueryByPropertyID(ctx context.Context, id uuid.UUID) ([]Floor, error) {
	flrs, err := c.store.QueryByPropertyID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("query: property id[%s]: %w", id.String(), err)
	}

	return flrs, nil
}

// ExecuteUnderTransaction constructs a new Core value that will use the
// specified transaction in any store related calls.
func (c *Core) ExecuteUnderTransaction(tx transaction.Transaction) (*Core, error) {
	store, err := c.store.ExecuteUnderTransaction(tx)
	if err != nil {
		return nil, err
	}

	c = &Core{
		store: store,
		log:   c.log,
	}

	return c, nil
}
