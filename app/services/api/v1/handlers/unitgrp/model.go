package unitgrp

import (
	"context"
	"github.com/google/uuid"
	"github.com/nhaancs/bhms/business/core/unit"
	"github.com/nhaancs/bhms/business/web/mid"
	"github.com/nhaancs/bhms/foundation/validate"
	"time"
)

// ==========================================================

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

// ==========================================================

type AppNewUnit struct {
	Name    string `json:"name" validate:"required"`
	BlockID string `json:"blockID" validate:"required"`
	FloorID string `json:"floorID" validate:"required"`
}

func (r AppNewUnit) Validate() error {
	if err := validate.Check(r); err != nil {
		return err
	}

	return nil
}

func toCoreNewUnit(ctx context.Context, a AppNewUnit) (unit.NewUnit, error) {
	blockID, err := uuid.Parse(a.BlockID)
	if err != nil {
		return unit.NewUnit{}, err
	}
	floorID, err := uuid.Parse(a.FloorID)
	if err != nil {
		return unit.NewUnit{}, err
	}
	return unit.NewUnit{
		ID:         uuid.New(),
		Name:       a.Name,
		PropertyID: mid.GetProperty(ctx).ID,
		BlockID:    blockID,
		FloorID:    floorID,
	}, nil
}

// ==========================================================

type AppUpdateUnit struct {
	Name *string `json:"name"`
}

func (r AppUpdateUnit) Validate() error {
	if err := validate.Check(r); err != nil {
		return err
	}

	return nil
}

func toCoreUpdateUnit(a AppUpdateUnit) (unit.UpdateUnit, error) {
	return unit.UpdateUnit{
		Name: a.Name,
	}, nil
}
