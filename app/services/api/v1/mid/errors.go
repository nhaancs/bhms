package mid

import (
	"context"
	"net/http"

	"github.com/nhaancs/bhms/app/services/api/v1/request"
	"github.com/nhaancs/bhms/business/web/auth"
	"github.com/nhaancs/bhms/foundation/logger"
	"github.com/nhaancs/bhms/foundation/validate"
	"github.com/nhaancs/bhms/foundation/web"
)

// Errors handles errors coming out of the call chain. It detects normal
// application errors which are used to respond to the client in a uniform way.
// Unexpected errors (status >= 500) are logged.
func Errors(log *logger.Logger) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			if err := handler(ctx, w, r); err != nil {
				log.Error(ctx, "message", "msg", err)

				ctx, span := web.AddSpan(ctx, "business.web.request.mid.error")
				span.RecordError(err)
				span.End()

				var er request.ErrorResponse
				var status int

				switch {
				case request.IsError(err):
					reqErr := request.GetError(err)

					if validate.IsFieldErrors(reqErr.Err) {
						fieldErrors := validate.GetFieldErrors(reqErr.Err)
						er = request.ErrorResponse{
							Error:  "data validation error",
							Fields: fieldErrors.Fields(),
						}
						status = reqErr.Status
						break
					}

					er = request.ErrorResponse{
						Error: reqErr.Error(),
					}
					status = reqErr.Status

				case auth.IsAuthError(err):
					er = request.ErrorResponse{
						Error: http.StatusText(http.StatusUnauthorized),
					}
					status = http.StatusUnauthorized

				default:
					er = request.ErrorResponse{
						Error: http.StatusText(http.StatusInternalServerError),
					}
					status = http.StatusInternalServerError
				}

				if err := web.Respond(ctx, w, er, status); err != nil {
					return err
				}

				// If we receive the shutdown err we need to return it
				// back to the base handler to shut down the service.
				if web.IsShutdown(err) {
					return err
				}
			}

			return nil
		}

		return h
	}

	return m
}
