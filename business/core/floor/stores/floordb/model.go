package floordb

import (
	"github.com/google/uuid"
	"github.com/nhaancs/bhms/business/core/floor"
	"time"
)

type dbFloor struct {
	ID         uuid.UUID `db:"id"`
	Name       string    `db:"name"`
	Status     string    `db:"status"`
	PropertyID uuid.UUID `db:"property_id"`
	BlockID    uuid.UUID `db:"block_id"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

func toDBFloor(c floor.Floor) dbFloor {
	return dbFloor{
		ID:         c.ID,
		Name:       c.Name,
		PropertyID: c.PropertyID,
		BlockID:    c.BlockID,
		Status:     c.Status.Name(),
		CreatedAt:  c.CreatedAt.UTC(),
		UpdatedAt:  c.UpdatedAt.UTC(),
	}
}

func toDBFloors(cs []floor.Floor) []dbFloor {
	result := make([]dbFloor, len(cs))
	for i := range cs {
		result[i] = toDBFloor(cs[i])
	}
	return result
}

func toCoreFloor(r dbFloor) (floor.Floor, error) {
	status, err := floor.ParseStatus(r.Status)
	if err != nil {
		return floor.Floor{}, err
	}
	return floor.Floor{
		ID:         r.ID,
		Name:       r.Name,
		PropertyID: r.PropertyID,
		BlockID:    r.BlockID,
		Status:     status,
		CreatedAt:  r.CreatedAt.In(time.Local),
		UpdatedAt:  r.UpdatedAt.In(time.Local),
	}, nil
}

func toCoreFloors(rs []dbFloor) ([]floor.Floor, error) {
	result := make([]floor.Floor, len(rs))
	for i := range rs {
		flr, err := toCoreFloor(rs[i])
		if err != nil {
			return nil, err
		}
		result[i] = flr
	}

	return result, nil
}
