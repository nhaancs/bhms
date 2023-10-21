// Package v1 manages the different versions of the API.
package v1

import (
	"github.com/nhaancs/bhms/app/services/api/v1/handlers/checkgrp"
	"github.com/nhaancs/bhms/app/services/api/v1/handlers/usergrp"
	"github.com/nhaancs/bhms/foundation/sms"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/nhaancs/bhms/app/services/api/v1/mid"
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
func APIMux(cfg APIMuxConfig, options ...func(opts *Options)) http.Handler {
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

	checkgrp.Routes(app, checkgrp.Config{
		Build: cfg.Build,
		DB:    cfg.DB,
	})

	usergrp.Routes(app, usergrp.Config{
		Log:   cfg.Log,
		Auth:  cfg.Auth,
		DB:    cfg.DB,
		KeyID: cfg.KeyID,
		SMS:   cfg.SMS,
	})

	return app
}
