package unitdb

import (
	"github.com/google/uuid"
	"github.com/nhaancs/bhms/business/core/unit"
	"time"
)

type dbUnit struct {
	ID         uuid.UUID `db:"id"`
	Name       string    `db:"name"`
	Status     string    `db:"status"`
	PropertyID uuid.UUID `db:"property_id"`
	BlockID    uuid.UUID `db:"block_id"`
	FloorID    uuid.UUID `db:"floor_id"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

func toDBUnit(c unit.Unit) dbUnit {
	return dbUnit{
		ID:         c.ID,
		Name:       c.Name,
		Status:     c.Status.Name(),
		PropertyID: c.PropertyID,
		BlockID:    c.BlockID,
		FloorID:    c.FloorID,
		CreatedAt:  c.CreatedAt.UTC(),
		UpdatedAt:  c.UpdatedAt.UTC(),
	}
}

func toDBUnits(cs []unit.Unit) []dbUnit {
	result := make([]dbUnit, len(cs))
	for i := range cs {
		result[i] = toDBUnit(cs[i])
	}
	return result
}

func toCoreUnit(r dbUnit) (unit.Unit, error) {
	status, err := unit.ParseStatus(r.Status)
	if err != nil {
		return unit.Unit{}, err
	}
	return unit.Unit{
		ID:         r.ID,
		Name:       r.Name,
		Status:     status,
		PropertyID: r.PropertyID,
		BlockID:    r.BlockID,
		FloorID:    r.FloorID,
		CreatedAt:  r.CreatedAt.In(time.Local),
		UpdatedAt:  r.UpdatedAt.In(time.Local),
	}, nil
}

func toCoreUnits(rs []dbUnit) ([]unit.Unit, error) {
	result := make([]unit.Unit, len(rs))
	for i := range rs {
		flr, err := toCoreUnit(rs[i])
		if err != nil {
			return nil, err
		}
		result[i] = flr
	}

	return result, nil
}
