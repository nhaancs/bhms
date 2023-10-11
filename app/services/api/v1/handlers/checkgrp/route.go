package checkgrp

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/nhaancs/bhms/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Build string
	DB    *sqlx.DB
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	hdl := New(cfg.Build, cfg.DB)
	app.HandleNoMiddleware(http.MethodGet, version, "/readiness", hdl.Readiness)
	app.HandleNoMiddleware(http.MethodGet, version, "/liveness", hdl.Liveness)
}
