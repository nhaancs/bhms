package usergrp

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/nhaancs/bhms/business/core/user"
	"github.com/nhaancs/bhms/business/core/user/stores/usercache"
	"github.com/nhaancs/bhms/business/core/user/stores/userdb"
	"github.com/nhaancs/bhms/business/web/auth"
	"github.com/nhaancs/bhms/foundation/logger"
	"github.com/nhaancs/bhms/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log   *logger.Logger
	Auth  *auth.Auth
	DB    *sqlx.DB
	KeyID string
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	usrCore := user.NewCore(cfg.Log, usercache.NewStore(cfg.Log, userdb.NewStore(cfg.Log, cfg.DB)))
	hdl := New(usrCore, cfg.Auth, cfg.KeyID)
	app.Handle(http.MethodPost, version, "/users/register", hdl.Register)
	app.Handle(http.MethodGet, version, "/users/token", hdl.Token)
}
