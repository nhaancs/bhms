package divisiongrp

import "github.com/nhaancs/bhms/business/core/division"

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
