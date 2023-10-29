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
	div *division.Core
}

// New constructs a handlers for route access.
func New(div *division.Core) *Handlers {
	return &Handlers{
		div: div,
	}
}

type AppDivision struct {
	ID       int    `json:"id"`
	ParentID int    `json:"parent_id"`
	Level    uint8  `json:"level"`
	Name     string `json:"name"`
}

func toAppDivision(d division.Division) AppDivision {
	return AppDivision{
		ID:       d.ID,
		ParentID: d.ParentID,
		Level:    d.Level,
		Name:     d.Name,
	}
}

func toAppDivisions(divs []division.Division) []AppDivision {
	result := make([]AppDivision, len(divs))

	for i := range divs {
		result[i] = toAppDivision(divs[i])
	}

	return result
}

func (h *Handlers) QueryProvinces(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	provinces, err := h.div.QueryProvinces(ctx)
	if err != nil {
		return fmt.Errorf("get province: %w", err)
	}

	return web.Respond(ctx, w, toAppDivisions(provinces), http.StatusOK)
}

func (h *Handlers) QueryByParentID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	parentID, err := strconv.ParseInt(web.Param(r, "parent_id"), 10, 32)
	if err != nil {
		return response.NewError(ErrInvalidID, http.StatusBadRequest)
	}

	divs, err := h.div.QueryByParentID(ctx, int(parentID))
	if err != nil {
		return fmt.Errorf("query by parrent id: %w", err)
	}

	return web.Respond(ctx, w, toAppDivisions(divs), http.StatusOK)
}
