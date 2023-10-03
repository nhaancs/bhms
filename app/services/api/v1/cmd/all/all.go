// Package all binds all the routes into the specified app.
package all

import (
	v1 "github.com/nhaancs/bhms/app/services/api/v1"
	"github.com/nhaancs/bhms/app/services/api/v1/handlers/checkgrp"
	"github.com/nhaancs/bhms/app/services/api/v1/handlers/usergrp"
	"github.com/nhaancs/bhms/foundation/web"
)

// Routes constructs the add value which provides the implementation of
// of RouteAdder for specifying what routes to bind to this instance.
func Routes() add {
	return add{}
}

type add struct{}

// Add implements the RouterAdder interface.
func (add) Add(app *web.App, cfg v1.APIMuxConfig) {
	checkgrp.Routes(app, checkgrp.Config{
		UsingWeaver: cfg.UsingWeaver,
		Build:       cfg.Build,
		DB:          cfg.DB,
	})

	usergrp.Routes(app, usergrp.Config{
		Log:  cfg.Log,
		Auth: cfg.Auth,
		DB:   cfg.DB,
	})
}
