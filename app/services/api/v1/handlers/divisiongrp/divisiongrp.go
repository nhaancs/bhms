package divisiongrp

import (
	"context"
	"fmt"
	"github.com/nhaancs/bhms/business/core/division"
	"github.com/nhaancs/bhms/foundation/web"
	"net/http"
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

func toAppDivision(d division.Divison) AppDivision {
	return AppDivision{
		ID:       d.ID,
		ParentID: d.ParentID,
		Level:    d.Level,
		Name:     d.Name,
	}
}

func toAppDivisions(divs []division.Divison) []AppDivision {
	result := make([]AppDivision, len(divs))

	for i := range divs {
		result[i] = toAppDivision(divs[i])
	}

	return result
}

func (h *Handlers) GetProvinces(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	provinces, err := h.div.GetProvinces(ctx)
	if err != nil {
		return fmt.Errorf("get province: %w", err)
	}

	return web.Respond(ctx, w, toAppDivisions(provinces), http.StatusCreated)
}
