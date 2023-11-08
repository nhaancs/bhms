package propertygrp

import (
	"github.com/google/uuid"
	"github.com/nhaancs/bhms/business/core/property"
	"github.com/nhaancs/bhms/foundation/validate"
	"time"
)

// ==========================================================

type AppProperty struct {
	ID              string     `json:"id"`
	ManagerID       string     `json:"managerID"`
	Name            string     `json:"name"`
	AddressLevel1ID uint32     `json:"addressLevel1ID"`
	AddressLevel2ID uint32     `json:"addressLevel2ID"`
	AddressLevel3ID uint32     `json:"addressLevel3ID"`
	Street          string     `json:"street"`
	Status          string     `json:"status"`
	CreatedAt       string     `json:"createdAt"`
	UpdatedAt       string     `json:"updatedAt"`
	Blocks          []AppBlock `json:"blocks"`
}

type AppBlock struct {
	ID         string     `json:"id"`
	PropertyID string     `json:"propertyID"`
	Name       string     `json:"name"`
	Floors     []AppFloor `json:"floors"`
	CreatedAt  string     `json:"createdAt"`
	UpdatedAt  string     `json:"updatedAt"`
}
type AppFloor struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	PropertyID string    `json:"propertyID"`
	BlockID    string    `json:"blockID"`
	Units      []AppUnit `json:"units"`
	CreatedAt  string    `json:"createdAt"`
	UpdatedAt  string    `json:"updatedAt"`
}
type AppUnit struct {
	ID         string `json:"id"`
	PropertyID string `json:"propertyID"`
	BlockID    string `json:"blockID"`
	FloorID    string `json:"floorID"`
	Name       string `json:"name"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
}

func toAppProperty(c property.Property) AppProperty {
	return AppProperty{
		ID:              c.ID.String(),
		ManagerID:       c.ManagerID.String(),
		AddressLevel1ID: c.AddressLevel1ID,
		AddressLevel2ID: c.AddressLevel2ID,
		AddressLevel3ID: c.AddressLevel3ID,
		Street:          c.Street,
		Status:          c.Status.Name(),
		CreatedAt:       c.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       c.UpdatedAt.Format(time.RFC3339),
	}
}
func toAppProperties(cs []property.Property) []AppProperty {
	result := make([]AppProperty, len(cs))
	for i := range cs {
		result[i] = toAppProperty(cs[i])
	}
	return result
}

// ===============================================================

type AppNewProperty struct {
	ManagerID       uuid.UUID     `json:"-"`
	Name            string        `json:"name" validate:"required"`
	AddressLevel1ID uint32        `json:"addressLevel1ID" validate:"required,min=1"`
	AddressLevel2ID uint32        `json:"addressLevel2ID" validate:"required,min=1"`
	AddressLevel3ID uint32        `json:"addressLevel3ID" validate:"required,min=1"`
	Street          string        `json:"street" validate:"required"`
	Status          string        `json:"status" validate:"required"`
	Blocks          []AppNewBlock `json:"blocks"`
}

type AppNewBlock struct {
	Name   string        `json:"name"`
	Floors []AppNewFloor `json:"floors"`
}
type AppNewFloor struct {
	Name  string       `json:"name"`
	Units []AppNewUnit `json:"units"`
}
type AppNewUnit struct {
	Name string `json:"name"`
}

func toCoreNewProperty(a AppNewProperty) (property.NewProperty, error) {
	return property.NewProperty{
		ManagerID:       a.ManagerID,
		Name:            a.Name,
		AddressLevel1ID: a.AddressLevel1ID,
		AddressLevel2ID: a.AddressLevel2ID,
		AddressLevel3ID: a.AddressLevel3ID,
		Street:          a.Street,
		Status:          property.StatusCreated,
	}, nil
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
	AddressLevel1ID *uint32 `json:"addressLevel1ID" validate:"required,min=1"`
	AddressLevel2ID *uint32 `json:"addressLevel2ID" validate:"required,min=1"`
	AddressLevel3ID *uint32 `json:"addressLevel3ID" validate:"required,min=1"`
	Street          *string `json:"street" validate:"required"`
}

func toCoreUpdateProperty(a AppUpdateProperty) (property.UpdateProperty, error) {
	return property.UpdateProperty{
		Name:            a.Name,
		AddressLevel1ID: a.AddressLevel1ID,
		AddressLevel2ID: a.AddressLevel2ID,
		AddressLevel3ID: a.AddressLevel3ID,
		Street:          a.Street,
	}, nil
}

// Validate checks the data in the model is considered clean.
func (r AppUpdateProperty) Validate() error {
	if err := validate.Check(r); err != nil {
		return err
	}

	return nil
}
