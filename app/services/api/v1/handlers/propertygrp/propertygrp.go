package propertygrp

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
	"github.com/nhaancs/bhms/business/web/auth"
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

func New(property *property.Core) *Handlers {
	return &Handlers{
		property: property,
	}
}

func (h *Handlers) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	h, err := h.executeUnderTransaction(ctx)
	if err != nil {
		return err
	}

	var app AppNewProperty
	if err := web.Decode(r, &app); err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}
	app.ManagerID = auth.GetUserID(ctx)

	c, err := toCoreNewProperty(app)
	if err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	prprty, err := h.property.Create(ctx, c)
	if err != nil {
		if errors.Is(err, property.ErrLimitExceeded) {
			return response.NewError(err, http.StatusForbidden)
		}
		return fmt.Errorf("create: prprty[%+v]: %w", app, err)
	}

	return web.Respond(ctx, w, toAppProperty(prprty), http.StatusCreated)
}

func (h *Handlers) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	prprtyIDStr := web.Param(r, "property_id")
	prprtyID, err := uuid.Parse(prprtyIDStr)
	if err != nil {
		return response.NewError(errors.New("invalid property id"), http.StatusBadRequest)
	}
	prprty, err := h.property.QueryByID(ctx, prprtyID)
	if err != nil {
		switch {
		case errors.Is(err, property.ErrNotFound):
			return response.NewError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querybyid: ID[%s]: %w", prprtyID, err)
		}
	}

	if prprty.ManagerID != auth.GetUserID(ctx) {
		return response.NewError(errors.New("no permission"), http.StatusForbidden)
	}

	var app AppUpdateProperty
	if err := web.Decode(r, &app); err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	c, err := toCoreUpdateProperty(app)
	if err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	prprty, err = h.property.Update(ctx, prprty, c)
	if err != nil {
		return fmt.Errorf("update: usr[%+v]: %w", app, err)
	}

	return web.Respond(ctx, w, toAppProperty(prprty), http.StatusCreated)
}

func (h *Handlers) Query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	managerID := auth.GetUserID(ctx)
	prprties, err := h.property.QueryByManagerID(ctx, managerID)
	if err != nil {
		switch {
		case errors.Is(err, property.ErrNotFound):
			return response.NewError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querybymanagerid: managerID[%s]: %w", managerID, err)
		}
	}

	return web.Respond(ctx, w, toAppProperties(prprties), http.StatusCreated)
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
