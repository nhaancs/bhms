package propertydb

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/nhaancs/bhms/business/core/property"
	"time"
)

type dbProperty struct {
	ID              uuid.UUID `db:"id"`
	ManagerID       uuid.UUID `db:"manager_id"`
	Name            string    `db:"name"`
	AddressLevel1ID uint32    `db:"address_level_1_id"`
	AddressLevel2ID uint32    `db:"address_level_2_id"`
	AddressLevel3ID uint32    `db:"address_level_3_id"`
	Street          string    `db:"street"`
	Status          string    `db:"status"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}

func toDBProperty(c property.Property) dbProperty {
	return dbProperty{
		ID:              c.ID,
		ManagerID:       c.ManagerID,
		Name:            c.Name,
		AddressLevel1ID: c.AddressLevel1ID,
		AddressLevel2ID: c.AddressLevel2ID,
		AddressLevel3ID: c.AddressLevel3ID,
		Street:          c.Street,
		Status:          c.Status.Name(),
		CreatedAt:       c.CreatedAt.UTC(),
		UpdatedAt:       c.UpdatedAt.UTC(),
	}
}

func toCoreProperty(r dbProperty) (property.Property, error) {
	status, err := property.ParseStatus(r.Status)
	if err != nil {
		return property.Property{}, fmt.Errorf("parse status: %w", err)
	}

	return property.Property{
		ID:              r.ID,
		ManagerID:       r.ManagerID,
		Name:            r.Name,
		AddressLevel1ID: r.AddressLevel1ID,
		AddressLevel2ID: r.AddressLevel2ID,
		AddressLevel3ID: r.AddressLevel3ID,
		Street:          r.Street,
		Status:          status,
		CreatedAt:       r.CreatedAt.In(time.Local),
		UpdatedAt:       r.UpdatedAt.In(time.Local),
	}, nil
}

func toCoreProperties(rs []dbProperty) ([]property.Property, error) {
	prprties := make([]property.Property, len(rs))
	var err error
	for i, r := range rs {
		prprties[i], err = toCoreProperty(r)
		if err != nil {
			return nil, err
		}
	}
	return prprties, nil
}
