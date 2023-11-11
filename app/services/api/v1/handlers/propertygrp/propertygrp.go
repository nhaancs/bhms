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

	if len(app.Blocks) > 5 {
		return response.NewError(errors.New("maximum number of blocks exceeded"), http.StatusBadRequest)
	}

	var (
		crNwPrprty     = toCoreNewProperty(app)
		crNwBlcks      = make([]block.NewBlock, len(app.Blocks))
		crNwFlrs       []floor.NewFloor
		appNwUntsMap   = make(map[[2]uuid.UUID][]AppNewUnit) // for convert O(n^3) to O(n^2)
		appNwUntsCount = 0
	)
	for i, appBlck := range app.Blocks {
		if len(appBlck.Floors) > 10 {
			return response.NewError(errors.New("maximum number of units of a floor exceeded"), http.StatusBadRequest)
		}

		crNwBlck := toCoreNewBlock(appBlck, crNwPrprty.ID)
		crNwBlcks[i] = crNwBlck

		for _, appFlr := range appBlck.Floors {
			{
				maxUnitsPerFloor := 200
				if len(appBlck.Floors) > 1 {
					maxUnitsPerFloor = 20
				}
				if len(appFlr.Units) > maxUnitsPerFloor {
					return response.NewError(errors.New("maximum number of floors of a block exceeded"), http.StatusBadRequest)
				}
			}

			crNwFlr := toCoreNewFloor(appFlr, crNwPrprty.ID, crNwBlck.ID)
			crNwFlrs = append(crNwFlrs, crNwFlr)

			appNwUntsMap[[2]uuid.UUID{crNwBlck.ID, crNwFlr.ID}] = appFlr.Units
			appNwUntsCount += len(appFlr.Units)
		}
	}

	crNwUnts := make([]unit.NewUnit, appNwUntsCount)
	for ids, appUnts := range appNwUntsMap {
		for i := range appUnts {
			crNwUnts[i] = toCoreNewUnit(appUnts[i], crNwPrprty.ID, ids[0], ids[1])
		}
	}

	prprty, err := h.property.Create(ctx, crNwPrprty)
	if err != nil {
		if errors.Is(err, property.ErrLimitExceeded) {
			return response.NewError(err, http.StatusBadRequest)
		}
		return fmt.Errorf("create: prprty[%+v]: %w", app, err)
	}

	blcks, err := h.block.BatchCreate(ctx, crNwBlcks)
	if err != nil {
		return fmt.Errorf("batch create blocks: prprty[%+v]: %w", app, err)
	}

	flrs, err := h.floor.BatchCreate(ctx, crNwFlrs)
	if err != nil {
		return fmt.Errorf("batch create floors: prprty[%+v]: %w", app, err)
	}

	unts, err := h.unit.BatchCreate(ctx, crNwUnts)
	if err != nil {
		return fmt.Errorf("batch create units: prprty[%+v]: %w", app, err)
	}

	return web.Respond(ctx, w, toAppPropertyDetail(prprty, blcks, flrs, unts), http.StatusCreated)
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

	appPrprtyDtls := make([]AppPropertyDetail, len(prprties))
	for i := range prprties {
		blcks, err := h.block.QueryByPropertyID(ctx, prprties[i].ID)
		if err != nil {
			return err
		}
		flrs, err := h.floor.QueryByPropertyID(ctx, prprties[i].ID)
		if err != nil {
			return err
		}
		unts, err := h.unit.QueryByPropertyID(ctx, prprties[i].ID)
		if err != nil {
			return err
		}

		appPrprtyDtls[i] = toAppPropertyDetail(prprties[i], blcks, flrs, unts)
	}

	return web.Respond(ctx, w, appPrprtyDtls, http.StatusCreated)
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
