package divisiongrp

import "github.com/nhaancs/bhms/business/core/division"

type AppDivision struct {
	ID       int    `json:"id"`
	ParentID int    `json:"parent_id"`
	Level    uint8  `json:"level"`
	Name     string `json:"name"`
}

func toAppDivision(c division.Division) AppDivision {
	return AppDivision{
		ID:       c.ID,
		ParentID: c.ParentID,
		Level:    c.Level,
		Name:     c.Name,
	}
}

func toAppDivisions(cs []division.Division) []AppDivision {
	result := make([]AppDivision, len(cs))
	for i := range cs {
		result[i] = toAppDivision(cs[i])
	}

	return result
}
