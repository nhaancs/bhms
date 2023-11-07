package propertygrp

import (
	"github.com/google/uuid"
	"github.com/nhaancs/bhms/business/core/property"
	"github.com/nhaancs/bhms/foundation/validate"
	"time"
)

// ==========================================================

type AppProperty struct {
	ID              string `json:"id"`
	ManagerID       string `json:"manager_id"`
	Name            string `json:"name"`
	AddressLevel1ID uint32 `json:"address_level_1_id"`
	AddressLevel2ID uint32 `json:"address_level_2_id"`
	AddressLevel3ID uint32 `json:"address_level_3_id"`
	Street          string `json:"street"`
	Status          string `json:"status"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

func toAppProperty(e property.Property) AppProperty {
	return AppProperty{
		ID:              e.ID.String(),
		ManagerID:       e.ManagerID.String(),
		AddressLevel1ID: e.AddressLevel1ID,
		AddressLevel2ID: e.AddressLevel2ID,
		AddressLevel3ID: e.AddressLevel3ID,
		Street:          e.Street,
		Status:          e.Status.Name(),
		CreatedAt:       e.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       e.UpdatedAt.Format(time.RFC3339),
	}
}
func toAppProperties(es []property.Property) []AppProperty {
	result := make([]AppProperty, len(es))
	for i := range es {
		result[i] = toAppProperty(es[i])
	}
	return result
}

// ===============================================================

type AppNewProperty struct {
	ManagerID       uuid.UUID `json:"-"`
	Name            string    `json:"name" validate:"required"`
	AddressLevel1ID uint32    `json:"address_level_1_id" validate:"required,min=1"`
	AddressLevel2ID uint32    `json:"address_level_2_id" validate:"required,min=1"`
	AddressLevel3ID uint32    `json:"address_level_3_id" validate:"required,min=1"`
	Street          string    `json:"street" validate:"required"`
	Status          string    `json:"status" validate:"required"`
}

func toCoreNewProperty(a AppNewProperty) (property.NewProperty, error) {
	prprty := property.NewProperty{
		ManagerID:       a.ManagerID,
		Name:            a.Name,
		AddressLevel1ID: a.AddressLevel1ID,
		AddressLevel2ID: a.AddressLevel2ID,
		AddressLevel3ID: a.AddressLevel3ID,
		Street:          a.Street,
		Status:          property.StatusCreated,
	}

	return prprty, nil
}

// Validate checks the data in the model is considered clean.
func (r AppNewProperty) Validate() error {
	if err := validate.Check(r); err != nil {
		return err
	}

	return nil
}

// ===============================================================

type AppUpdateProperty struct {
	Name            *string `json:"name" validate:"required"`
	AddressLevel1ID *uint32 `json:"address_level_1_id" validate:"required,min=1"`
	AddressLevel2ID *uint32 `json:"address_level_2_id" validate:"required,min=1"`
	AddressLevel3ID *uint32 `json:"address_level_3_id" validate:"required,min=1"`
	Street          *string `json:"street" validate:"required"`
}

func toCoreUpdateProperty(a AppUpdateProperty) (property.UpdateProperty, error) {
	prprty := property.UpdateProperty{
		Name:            a.Name,
		AddressLevel1ID: a.AddressLevel1ID,
		AddressLevel2ID: a.AddressLevel2ID,
		AddressLevel3ID: a.AddressLevel3ID,
		Street:          a.Street,
	}

	return prprty, nil
}

// Validate checks the data in the model is considered clean.
func (r AppUpdateProperty) Validate() error {
	if err := validate.Check(r); err != nil {
		return err
	}

	return nil
}
