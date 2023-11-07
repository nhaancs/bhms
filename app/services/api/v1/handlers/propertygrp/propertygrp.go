package propertygrp

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/nhaancs/bhms/business/core/property"
	"github.com/nhaancs/bhms/business/web/auth"
	"github.com/nhaancs/bhms/business/web/response"
	"github.com/nhaancs/bhms/foundation/web"
	"net/http"
)

type Handlers struct {
	property *property.Core
}

func New(
	property *property.Core,
) *Handlers {
	return &Handlers{
		property: property,
	}
}

func (h *Handlers) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app AppNewProperty
	if err := web.Decode(r, &app); err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}
	app.ManagerID = auth.GetUserID(ctx)

	e, err := toCoreNewProperty(app)
	if err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	prprty, err := h.property.Create(ctx, e)
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

	e, err := toCoreUpdateProperty(app)
	if err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	prprty, err = h.property.Update(ctx, prprty, e)
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
