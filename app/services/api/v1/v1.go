// Package v1 manages the different versions of the API.
package v1

import (
	"github.com/nhaancs/bhms/app/services/api/v1/handlers/checkgrp"
	"github.com/nhaancs/bhms/app/services/api/v1/handlers/divisiongrp"
	"github.com/nhaancs/bhms/app/services/api/v1/handlers/propertygrp"
	"github.com/nhaancs/bhms/app/services/api/v1/handlers/usergrp"
	"github.com/nhaancs/bhms/business/core/division"
	"github.com/nhaancs/bhms/business/core/division/stores/divisionjson"
	"github.com/nhaancs/bhms/business/core/property"
	"github.com/nhaancs/bhms/business/core/property/stores/propertycache"
	"github.com/nhaancs/bhms/business/core/property/stores/propertydb"
	"github.com/nhaancs/bhms/business/core/user"
	"github.com/nhaancs/bhms/business/core/user/stores/usercache"
	"github.com/nhaancs/bhms/business/core/user/stores/userdb"
	"github.com/nhaancs/bhms/business/web/mid"
	"github.com/nhaancs/bhms/foundation/sms"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/nhaancs/bhms/business/web/auth"
	"github.com/nhaancs/bhms/foundation/logger"
	"github.com/nhaancs/bhms/foundation/web"
	"go.opentelemetry.io/otel/trace"
)

// Options represent optional parameters.
type Options struct {
	corsOrigin string
}

// WithCORS provides configuration options for CORS.
func WithCORS(origin string) func(opts *Options) {
	return func(opts *Options) {
		opts.corsOrigin = origin
	}
}

// APIMuxConfig contains all the mandatory systems required by handlers.
type APIMuxConfig struct {
	Build    string
	Shutdown chan os.Signal
	Log      *logger.Logger
	Auth     *auth.Auth
	DB       *sqlx.DB
	Tracer   trace.Tracer
	KeyID    string
	SMS      *sms.SMS
}

// APIMux constructs a http.Handler with all application routes defined.
func APIMux(cfg APIMuxConfig, options ...func(opts *Options)) (http.Handler, error) {
	const version = "v1"
	var opts Options
	for _, option := range options {
		option(&opts)
	}

	app := web.NewApp(
		cfg.Shutdown,
		cfg.Tracer,
		mid.Logger(cfg.Log),
		mid.Errors(cfg.Log),
		mid.Metrics(),
		mid.Panics(),
	)

	if opts.corsOrigin != "" {
		app.EnableCORS(mid.Cors(opts.corsOrigin))
	}

	// -------------------------------------------------------------------------
	// Check routes
	checkHdl := checkgrp.New(cfg.Build, cfg.DB)
	app.HandleNoMiddleware(http.MethodGet, version, "/readiness", checkHdl.Readiness)
	app.HandleNoMiddleware(http.MethodGet, version, "/liveness", checkHdl.Liveness)

	// -------------------------------------------------------------------------
	// User routes
	usrCore := user.NewCore(cfg.Log, usercache.NewStore(cfg.Log, userdb.NewStore(cfg.Log, cfg.DB)))
	usrHdl := usergrp.New(usrCore, cfg.Auth, cfg.KeyID, cfg.SMS)
	app.Handle(http.MethodPost, version, "/users/register", usrHdl.Register)
	app.Handle(http.MethodPost, version, "/users/verify-otp", usrHdl.VerifyOTP)
	app.Handle(http.MethodGet, version, "/users/token", usrHdl.Token)
	app.Handle(http.MethodPut, version, "/users", usrHdl.Token, mid.Authenticate(cfg.Auth))

	// -------------------------------------------------------------------------
	// Division routes
	divStore, err := divisionjson.NewStore(cfg.Log)
	if err != nil {
		return nil, err
	}
	divCore := division.NewCore(cfg.Log, divStore)
	divHdl := divisiongrp.New(divCore)
	app.Handle(http.MethodGet, version, "/divisions/provinces", divHdl.QueryProvinces, mid.Authenticate(cfg.Auth))
	app.Handle(http.MethodGet, version, "/divisions/children/:parent_id", divHdl.QueryByParentID, mid.Authenticate(cfg.Auth))

	// -------------------------------------------------------------------------
	// Property routes

	propertyStore := propertycache.NewStore(cfg.Log, propertydb.NewStore(cfg.Log, cfg.DB))
	if err != nil {
		return nil, err
	}
	propertyCore := property.NewCore(cfg.Log, propertyStore)
	propertyHdl := propertygrp.New(propertyCore)
	app.Handle(http.MethodGet, version, "/properties", propertyHdl.Create, mid.Authenticate(cfg.Auth))
	app.Handle(http.MethodPut, version, "/properties/:property_id", propertyHdl.Update, mid.Authenticate(cfg.Auth))

	return app, nil
}
