package propertygrp

import (
	"context"
	"errors"
	"fmt"
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

// TODO: update, list by manager id

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
