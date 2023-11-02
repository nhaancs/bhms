package tests

import (
	"github.com/nhaancs/bhms/app/services/api/v1/handlers/usergrp"
	"github.com/nhaancs/bhms/business/core/user"
	"time"
)

type tableData struct {
	name       string
	url        string
	token      string
	method     string
	statusCode int
	model      any
	resp       any
	expResp    any
	cmpFunc    func(x interface{}, y interface{}) string
}

type seedData struct {
}

// =============================================================================

func toAppUser(usr user.User) usergrp.AppUser {
	roles := make([]string, len(usr.Roles))
	for i, role := range usr.Roles {
		roles[i] = role.Name()
	}

	return usergrp.AppUser{
		ID:           usr.ID.String(),
		FirstName:    usr.FirstName,
		LastName:     usr.LastName,
		Phone:        usr.Phone,
		Status:       usr.Status.Name(),
		Roles:        roles,
		PasswordHash: nil, // This field is not marshalled.
		CreatedAt:    usr.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    usr.UpdatedAt.Format(time.RFC3339),
	}
}

func toAppUsers(users []user.User) []usergrp.AppUser {
	items := make([]usergrp.AppUser, len(users))
	for i, usr := range users {
		items[i] = toAppUser(usr)
	}

	return items
}

func toAppUserPtr(usr user.User) *usergrp.AppUser {
	appUsr := toAppUser(usr)
	return &appUsr
}

// =============================================================================
