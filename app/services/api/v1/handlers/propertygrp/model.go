package propertygrp

import (
	"github.com/google/uuid"
	"github.com/nhaancs/bhms/business/core/block"
	"github.com/nhaancs/bhms/business/core/floor"
	"github.com/nhaancs/bhms/business/core/property"
	"github.com/nhaancs/bhms/business/core/unit"
	"github.com/nhaancs/bhms/foundation/validate"
	"time"
)

// ==========================================================

type AppProperty struct {
	ID              string `json:"id"`
	ManagerID       string `json:"managerID"`
	Name            string `json:"name"`
	AddressLevel1ID uint32 `json:"addressLevel1ID"`
	AddressLevel2ID uint32 `json:"addressLevel2ID"`
	AddressLevel3ID uint32 `json:"addressLevel3ID"`
	Street          string `json:"street"`
	Status          string `json:"status"`
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updatedAt"`
}

func toAppProperty(p property.Property) AppProperty {
	return AppProperty{
		ID:              p.ID.String(),
		ManagerID:       p.ManagerID.String(),
		AddressLevel1ID: p.AddressLevel1ID,
		AddressLevel2ID: p.AddressLevel2ID,
		AddressLevel3ID: p.AddressLevel3ID,
		Street:          p.Street,
		Status:          p.Status.Name(),
		CreatedAt:       p.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       p.UpdatedAt.Format(time.RFC3339),
	}
}

// ==========================================================

type AppPropertyDetail struct {
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
	Status     string     `json:"status"`
	CreatedAt  string     `json:"createdAt"`
	UpdatedAt  string     `json:"updatedAt"`
}
type AppFloor struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	PropertyID string    `json:"propertyID"`
	BlockID    string    `json:"blockID"`
	Units      []AppUnit `json:"units"`
	Status     string    `json:"status"`
	CreatedAt  string    `json:"createdAt"`
	UpdatedAt  string    `json:"updatedAt"`
}
type AppUnit struct {
	ID         string `json:"id"`
	PropertyID string `json:"propertyID"`
	BlockID    string `json:"blockID"`
	FloorID    string `json:"floorID"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
}

func toAppPropertyDetail(prprty property.Property, blcks []block.Block, flrs []floor.Floor, unts []unit.Unit) AppPropertyDetail {
	appUntsMap := make(map[[2]uuid.UUID][]AppUnit) // [blockID, floorID] => units
	for i := range unts {
		key := [2]uuid.UUID{unts[i].BlockID, unts[i].FloorID}
		appUntsMap[key] = append(appUntsMap[key], toAppUnit(unts[i]))
	}

	appFlrsMap := make(map[uuid.UUID][]AppFloor) // [blockID] => floor
	for j := range flrs {
		flrKey := flrs[j].BlockID
		untKey := [2]uuid.UUID{flrs[j].BlockID, flrs[j].ID}
		appFlrsMap[flrKey] = append(appFlrsMap[flrKey], toAppFloor(flrs[j], appUntsMap[untKey]))
	}

	appBlcks := make([]AppBlock, len(blcks))
	for k := range blcks {
		appBlcks[k] = toAppBlock(blcks[k], appFlrsMap[blcks[k].ID])
	}

	return AppPropertyDetail{
		ID:              prprty.ID.String(),
		ManagerID:       prprty.ManagerID.String(),
		AddressLevel1ID: prprty.AddressLevel1ID,
		AddressLevel2ID: prprty.AddressLevel2ID,
		AddressLevel3ID: prprty.AddressLevel3ID,
		Street:          prprty.Street,
		Status:          prprty.Status.Name(),
		CreatedAt:       prprty.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       prprty.UpdatedAt.Format(time.RFC3339),
		Blocks:          appBlcks,
	}
}

func toAppUnit(u unit.Unit) AppUnit {
	return AppUnit{
		ID:         u.ID.String(),
		PropertyID: u.PropertyID.String(),
		BlockID:    u.BlockID.String(),
		FloorID:    u.FloorID.String(),
		Name:       u.Name,
		Status:     u.Status.Name(),
		CreatedAt:  u.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  u.UpdatedAt.Format(time.RFC3339),
	}
}

func toAppFloor(f floor.Floor, appUnts []AppUnit) AppFloor {
	return AppFloor{
		ID:         f.ID.String(),
		PropertyID: f.PropertyID.String(),
		BlockID:    f.BlockID.String(),
		Name:       f.Name,
		Status:     f.Status.Name(),
		CreatedAt:  f.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  f.UpdatedAt.Format(time.RFC3339),
		Units:      appUnts,
	}
}

func toAppBlock(b block.Block, appFlrs []AppFloor) AppBlock {
	return AppBlock{
		ID:         b.ID.String(),
		PropertyID: b.PropertyID.String(),
		Name:       b.Name,
		Status:     b.Status.Name(),
		CreatedAt:  b.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  b.UpdatedAt.Format(time.RFC3339),
		Floors:     appFlrs,
	}
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
	Blocks          []AppNewBlock `json:"blocks" validate:"required"`
}

type AppNewBlock struct {
	Name   string        `json:"name" validate:"required"`
	Floors []AppNewFloor `json:"floors" validate:"required"`
}
type AppNewFloor struct {
	Name  string       `json:"name" validate:"required"`
	Units []AppNewUnit `json:"units" validate:"required"`
}
type AppNewUnit struct {
	Name string `json:"name" validate:"required"`
}

func toCoreNewProperty(a AppNewProperty) property.NewProperty {
	return property.NewProperty{
		ID:              uuid.New(),
		ManagerID:       a.ManagerID,
		Name:            a.Name,
		AddressLevel1ID: a.AddressLevel1ID,
		AddressLevel2ID: a.AddressLevel2ID,
		AddressLevel3ID: a.AddressLevel3ID,
		Street:          a.Street,
		Status:          property.StatusCreated,
	}
}

// Validate checks the data in the model is considered clean.
func (r AppNewProperty) Validate() error {
	if err := validate.Check(r); err != nil {
		return err
	}

	return nil
}

func toCoreNewBlock(a AppNewBlock, propertyID uuid.UUID) block.NewBlock {
	return block.NewBlock{
		ID:         uuid.New(),
		Name:       a.Name,
		PropertyID: propertyID,
	}
}
func toCoreNewFloor(a AppNewFloor, propertyID, blockID uuid.UUID) floor.NewFloor {
	return floor.NewFloor{
		ID:         uuid.New(),
		Name:       a.Name,
		PropertyID: propertyID,
		BlockID:    blockID,
	}
}

func toCoreNewUnit(a AppNewUnit, propertyID, blockID, floorID uuid.UUID) unit.NewUnit {
	return unit.NewUnit{
		ID:         uuid.New(),
		Name:       a.Name,
		PropertyID: propertyID,
		BlockID:    blockID,
		FloorID:    floorID,
	}
}

// ===============================================================

type AppUpdateProperty struct {
	Name            *string `json:"name" validate:"required"`
	AddressLevel1ID *uint32 `json:"addressLevel1ID" validate:"required,min=1"`
	AddressLevel2ID *uint32 `json:"addressLevel2ID" validate:"required,min=1"`
	AddressLevel3ID *uint32 `json:"addressLevel3ID" validate:"required,min=1"`
	Street          *string `json:"street" validate:"required"`
	Status          *string `json:"status" validate:"required"`
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
