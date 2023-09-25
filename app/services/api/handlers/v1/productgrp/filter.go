package productgrp

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/nhaancs/realworld/business/core/product"
	"github.com/nhaancs/realworld/foundation/validate"
)

func parseFilter(r *http.Request) (product.QueryFilter, error) {
	const (
		filterByProdID   = "product_id"
		filterByCost     = "cost"
		filterByQuantity = "quantity"
		filterByName     = "name"
	)

	values := r.URL.Query()

	var filter product.QueryFilter

	if productID := values.Get(filterByProdID); productID != "" {
		id, err := uuid.Parse(productID)
		if err != nil {
			return product.QueryFilter{}, validate.NewFieldsError(filterByProdID, err)
		}
		filter.WithProductID(id)
	}

	if cost := values.Get(filterByCost); cost != "" {
		cst, err := strconv.ParseFloat(cost, 64)
		if err != nil {
			return product.QueryFilter{}, validate.NewFieldsError(filterByCost, err)
		}
		filter.WithCost(cst)
	}

	if quantity := values.Get(filterByQuantity); quantity != "" {
		qua, err := strconv.ParseInt(quantity, 10, 64)
		if err != nil {
			return product.QueryFilter{}, validate.NewFieldsError(filterByQuantity, err)
		}
		filter.WithQuantity(int(qua))
	}

	if name := values.Get(filterByName); name != "" {
		filter.WithName(name)
	}

	if err := filter.Validate(); err != nil {
		return product.QueryFilter{}, err
	}

	return filter, nil
}
