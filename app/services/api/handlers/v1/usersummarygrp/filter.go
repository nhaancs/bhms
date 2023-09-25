package usersummarygrp

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/nhaancs/realworld/business/core/usersummary"
	"github.com/nhaancs/realworld/foundation/validate"
)

func parseFilter(r *http.Request) (usersummary.QueryFilter, error) {
	const (
		filterByUserID = "user_id"
		filterByName   = "name"
	)

	values := r.URL.Query()

	var filter usersummary.QueryFilter

	if userID := values.Get(filterByUserID); userID != "" {
		id, err := uuid.Parse(userID)
		if err != nil {
			return usersummary.QueryFilter{}, validate.NewFieldsError(filterByUserID, err)
		}
		filter.WithUserID(id)
	}

	if userName := values.Get(filterByName); userName != "" {
		filter.WithUserName(userName)
	}

	return filter, nil
}
