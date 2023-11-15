package unitgrp

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/nhaancs/bhms/business/core/block"
	"github.com/nhaancs/bhms/business/core/floor"
	"github.com/nhaancs/bhms/business/core/property"
	"github.com/nhaancs/bhms/business/core/unit"
	"github.com/nhaancs/bhms/business/data/transaction"
	"github.com/nhaancs/bhms/business/web/mid"
	"github.com/nhaancs/bhms/business/web/response"
	"github.com/nhaancs/bhms/foundation/web"
	"net/http"
)

type Handlers struct {
	property *property.Core
	block    *block.Core
	floor    *floor.Core
	unit     *unit.Core
}

func New(
	property *property.Core,
	block *block.Core,
	floor *floor.Core,
	unit *unit.Core,
) *Handlers {
	return &Handlers{
		property: property,
		block:    block,
		floor:    floor,
		unit:     unit,
	}
}

func (h *Handlers) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app AppNewUnit
	if err := web.Decode(r, &app); err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	c, err := toCoreNewUnit(ctx, app)
	if err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	prprty := mid.GetProperty(ctx)
	blcks, err := h.block.QueryByPropertyID(ctx, prprty.ID)
	if err != nil {
		switch {
		case errors.Is(err, block.ErrNotFound):
			return response.NewError(err, http.StatusNoContent)
		default:
			return fmt.Errorf("QueryByPropertyID: id[%s]: %w", prprty.ID, err)
		}
	}
	if len(blcks) >= 1000 {
		return response.NewError(errors.New("maximum number of units exceeded"), http.StatusBadRequest)
	}

	_, err = h.block.QueryByID(ctx, c.BlockID)
	if err != nil {
		switch {
		case errors.Is(err, block.ErrNotFound):
			return response.NewError(err, http.StatusNoContent)
		default:
			return fmt.Errorf("querybyid: block id[%s]: %w", c.BlockID, err)
		}
	}

	_, err = h.floor.QueryByID(ctx, c.FloorID)
	if err != nil {
		switch {
		case errors.Is(err, floor.ErrNotFound):
			return response.NewError(err, http.StatusNoContent)
		default:
			return fmt.Errorf("querybyid: floor id[%s]: %w", c.FloorID, err)
		}
	}

	unt, err := h.unit.Create(ctx, c)
	if err != nil {
		return fmt.Errorf("create: unit[%+v]: %w", app, err)
	}

	return web.Respond(ctx, w, toAppUnit(unt), http.StatusOK)
}

func (h *Handlers) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app AppUpdateUnit
	if err := web.Decode(r, &app); err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	c, err := toCoreUpdateUnit(app)
	if err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	unitID, err := uuid.Parse(web.Param(r, "unit_id"))
	if err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}
	unt, err := h.unit.QueryByID(ctx, unitID)
	if err != nil {
		switch {
		case errors.Is(err, unit.ErrNotFound):
			return response.NewError(err, http.StatusNoContent)
		default:
			return fmt.Errorf("querybyid: unit id[%s]: %w", unitID, err)
		}
	}

	unt, err = h.unit.Update(ctx, unt, c)
	if err != nil {
		return fmt.Errorf("update: unt[%+v]: %w", app, err)
	}

	return web.Respond(ctx, w, toAppUnit(unt), http.StatusOK)
}

func (h *Handlers) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	h, err := h.executeUnderTransaction(ctx)
	if err != nil {
		return err
	}

	unitID, err := uuid.Parse(web.Param(r, "unit_id"))
	if err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}
	unt, err := h.unit.QueryByID(ctx, unitID)
	if err != nil {
		switch {
		case errors.Is(err, unit.ErrNotFound):
			return response.NewError(err, http.StatusNoContent)
		default:
			return fmt.Errorf("querybyid: unit id[%s]: %w", unitID, err)
		}
	}
	unt, err = h.unit.Delete(ctx, unt)
	if err != nil {
		return fmt.Errorf("delete: unit id[%s]: %w", unt.ID.String(), err)
	}

	return web.Respond(ctx, w, toAppUnit(unt), http.StatusOK)
}

// executeUnderTransaction constructs a new Handlers value with the core apis
// using a store transaction that was created via middleware.
func (h *Handlers) executeUnderTransaction(ctx context.Context) (*Handlers, error) {
	if tx, ok := transaction.Get(ctx); ok {
		property, err := h.property.ExecuteUnderTransaction(tx)
		if err != nil {
			return nil, err
		}

		block, err := h.block.ExecuteUnderTransaction(tx)
		if err != nil {
			return nil, err
		}

		floor, err := h.floor.ExecuteUnderTransaction(tx)
		if err != nil {
			return nil, err
		}

		unit, err := h.unit.ExecuteUnderTransaction(tx)
		if err != nil {
			return nil, err
		}

		h = &Handlers{
			property: property,
			block:    block,
			floor:    floor,
			unit:     unit,
		}

		return h, nil
	}

	return h, nil
}
