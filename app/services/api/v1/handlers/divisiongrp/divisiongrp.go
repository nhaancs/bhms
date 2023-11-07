// Package divisiongrp ...
package divisiongrp

import (
	"context"
	"errors"
	"fmt"
	"github.com/nhaancs/bhms/business/core/division"
	"github.com/nhaancs/bhms/business/web/response"
	"github.com/nhaancs/bhms/foundation/web"
	"net/http"
	"strconv"
)

// Set of error variables for handling group errors.
var (
	ErrInvalidID = errors.New("invalid id")
)

// Handlers manages the set of user endpoints.
type Handlers struct {
	division *division.Core
}

// New constructs a handlers for route access.
func New(division *division.Core) *Handlers {
	return &Handlers{
		division: division,
	}
}

func (h *Handlers) QueryProvinces(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	prvncs, err := h.division.QueryProvinces(ctx)
	if err != nil {
		return fmt.Errorf("get province: %w", err)
	}

	return web.Respond(ctx, w, toAppDivisions(prvncs), http.StatusOK)
}

func (h *Handlers) QueryByParentID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	parentID, err := strconv.ParseInt(web.Param(r, "parent_id"), 10, 32)
	if err != nil {
		return response.NewError(ErrInvalidID, http.StatusBadRequest)
	}

	dvsn, err := h.division.QueryByParentID(ctx, int(parentID))
	if err != nil {
		return fmt.Errorf("query by parrent id: %w", err)
	}

	return web.Respond(ctx, w, toAppDivisions(dvsn), http.StatusOK)
}
