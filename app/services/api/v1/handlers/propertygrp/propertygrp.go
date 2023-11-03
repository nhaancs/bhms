package propertygrp

import (
	"github.com/nhaancs/bhms/business/core/property"
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

// TODO: create (limit 1 per user), update, list by manager id
