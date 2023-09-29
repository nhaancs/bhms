// Package v1 contains the full set of handler functions and routes
// supported by the v1 web api.
package v1

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/nhaancs/realworld/app/services/api/handlers/v1/checkgrp"
	"github.com/nhaancs/realworld/app/services/api/handlers/v1/usergrp"
	"github.com/nhaancs/realworld/business/core/event"
	"github.com/nhaancs/realworld/business/core/user"
	"github.com/nhaancs/realworld/business/core/user/stores/usercache"
	"github.com/nhaancs/realworld/business/core/user/stores/userdb"
	db "github.com/nhaancs/realworld/business/data/dbsql/pgx"
	"github.com/nhaancs/realworld/business/web/auth"
	"github.com/nhaancs/realworld/business/web/v1/mid"
	"github.com/nhaancs/realworld/foundation/logger"
	"github.com/nhaancs/realworld/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Build string
	Log   *logger.Logger
	Auth  *auth.Auth
	DB    *sqlx.DB
}

// Routes binds all the version 1 routes.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	envCore := event.NewCore(cfg.Log)
	usrCore := user.NewCore(cfg.Log, envCore, usercache.NewStore(cfg.Log, userdb.NewStore(cfg.Log, cfg.DB)))

	authen := mid.Authenticate(cfg.Auth)
	tran := mid.ExecuteInTransation(cfg.Log, db.NewBeginner(cfg.DB))

	// -------------------------------------------------------------------------

	cgh := checkgrp.New(cfg.Build, cfg.DB)

	app.HandleNoMiddleware(http.MethodGet, version, "/readiness", cgh.Readiness)
	app.HandleNoMiddleware(http.MethodGet, version, "/liveness", cgh.Liveness)

	// -------------------------------------------------------------------------

	ugh := usergrp.New(usrCore, cfg.Auth)

	app.Handle(http.MethodGet, version, "/users/token/:kid", ugh.Token)
	app.Handle(http.MethodGet, version, "/users/:user_id", ugh.QueryByID, authen)
	app.Handle(http.MethodPost, version, "/users", ugh.Create, authen)
	app.Handle(http.MethodPut, version, "/users/:user_id", ugh.Update, authen, tran)
}
