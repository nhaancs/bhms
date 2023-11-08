package blockdb

import (
	"github.com/google/uuid"
	"github.com/nhaancs/bhms/business/core/block"
	"time"
)

type dbBlock struct {
	ID         uuid.UUID `db:"id"`
	Name       string    `db:"name"`
	PropertyID uuid.UUID `db:"property_id"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

func toDBBlock(c block.Block) dbBlock {
	return dbBlock{
		ID:         c.ID,
		Name:       c.Name,
		PropertyID: c.PropertyID,
		CreatedAt:  c.CreatedAt.UTC(),
		UpdatedAt:  c.UpdatedAt.UTC(),
	}
}

func toDBBlocks(cs []block.Block) []dbBlock {
	result := make([]dbBlock, len(cs))
	for i := range cs {
		result[i] = toDBBlock(cs[i])
	}
	return result
}

func toCoreBlock(r dbBlock) (block.Block, error) {
	return block.Block{
		ID:         r.ID,
		Name:       r.Name,
		PropertyID: r.PropertyID,
		CreatedAt:  r.CreatedAt.In(time.Local),
		UpdatedAt:  r.UpdatedAt.In(time.Local),
	}, nil
}

func toCoreBlocks(rs []dbBlock) ([]block.Block, error) {
	result := make([]block.Block, len(rs))
	for i := range rs {
		blck, err := toCoreBlock(rs[i])
		if err != nil {
			return nil, err
		}
		result[i] = blck
	}

	return result, nil
}
