package division

import (
	"context"
	"errors"
	"fmt"
	"github.com/nhaancs/bhms/foundation/logger"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound = errors.New("division not found")
)

// =============================================================================

type Storer interface {
	QueryByID(ctx context.Context, divisionID int) (Divison, error)
	QueryByParentID(ctx context.Context, parentID int) ([]Divison, error)
	QueryLevel1s(ctx context.Context) ([]Divison, error)
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

func (c *Core) QueryByID(ctx context.Context, divisionID int) (Divison, error) {
	div, err := c.store.QueryByID(ctx, divisionID)
	if err != nil {
		return Divison{}, fmt.Errorf("query: divisionID[%d]: err: %w", divisionID, err)
	}

	return div, err
}

func (c *Core) QueryByParentID(ctx context.Context, parentID int) ([]Divison, error) {
	divs, err := c.store.QueryByParentID(ctx, parentID)
	if err != nil {
		return nil, fmt.Errorf("query: parentID[%d]: err: %w", parentID, err)
	}

	return divs, err
}
