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

func toDBProperty(e property.Property) dbProperty {
	return dbProperty{
		ID:              e.ID,
		ManagerID:       e.ManagerID,
		Name:            e.Name,
		AddressLevel1ID: e.AddressLevel1ID,
		AddressLevel2ID: e.AddressLevel2ID,
		AddressLevel3ID: e.AddressLevel3ID,
		Street:          e.Street,
		Status:          e.Status.Name(),
		CreatedAt:       e.CreatedAt.UTC(),
		UpdatedAt:       e.UpdatedAt.UTC(),
	}
}

func toCoreProperty(dbProperty dbProperty) (property.Property, error) {
	status, err := property.ParseStatus(dbProperty.Status)
	if err != nil {
		return property.Property{}, fmt.Errorf("parse status: %w", err)
	}

	prprty := property.Property{
		ID:              dbProperty.ID,
		ManagerID:       dbProperty.ManagerID,
		Name:            dbProperty.Name,
		AddressLevel1ID: dbProperty.AddressLevel1ID,
		AddressLevel2ID: dbProperty.AddressLevel2ID,
		AddressLevel3ID: dbProperty.AddressLevel3ID,
		Street:          dbProperty.Street,
		Status:          status,
		CreatedAt:       dbProperty.CreatedAt.In(time.Local),
		UpdatedAt:       dbProperty.UpdatedAt.In(time.Local),
	}

	return prprty, nil
}

func toCoreProperties(rows []dbProperty) ([]property.Property, error) {
	prprties := make([]property.Property, len(rows))
	var err error
	for i, dbPrprty := range rows {
		prprties[i], err = toCoreProperty(dbPrprty)
		if err != nil {
			return nil, err
		}
	}
	return prprties, nil
}
