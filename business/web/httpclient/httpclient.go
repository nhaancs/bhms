package httpclient

import (
	"fmt"
	"github.com/nhaancs/bhms/foundation/logger"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type options struct {
	logger  *logger.Logger
	logBody bool
	tracing bool
	metrics bool
	proxy   func(request *http.Request) (*url.URL, error)
}

type Option func(o *options)

func WithLogger(l *logger.Logger, body bool) Option {
	return func(o *options) {
		o.logger = l
		o.logBody = body
	}
}

func WithTracing() Option {
	return func(o *options) {
		o.tracing = true
	}
}

func WithMetrics() Option {
	return func(o *options) {
		o.metrics = true
	}
}

func WithProxy(proxy *url.URL) Option {
	return func(o *options) {
		if proxy != nil {
			o.proxy = http.ProxyURL(proxy)
		}
	}
}

type roundTripperFn func(req *http.Request) (*http.Response, error)

func (f roundTripperFn) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

var roundTripper http.RoundTripper = &http.Transport{
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext,
	ForceAttemptHTTP2:     true,
	MaxIdleConns:          1000,
	MaxIdleConnsPerHost:   100,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}

// New returns a HTTP client with logging, tracing, metric, and proxy support
func New(opts ...Option) *http.Client {
	o := new(options)
	for _, opt := range opts {
		opt(o)
	}

	rt := roundTripper
	if t, ok := rt.(*http.Transport); ok && o.proxy != nil {
		t = t.Clone()
		t.Proxy = o.proxy
		rt = t
	}

	if o.tracing {
		rt = otelhttp.NewTransport(rt)
	}

	if o.logger != nil {
		rt = logRoundTripper(rt, o.logger, o.logBody)
	}

	if o.metrics {
		rt = metricsRoundTripper(rt)
	}

	return &http.Client{
		Timeout:   30 * time.Second,
		Transport: rt,
	}
}

func logRoundTripper(rt http.RoundTripper, l *logger.Logger, body bool) http.RoundTripper {
	return roundTripperFn(func(req *http.Request) (*http.Response, error) {
		ctx := req.Context()
		start := time.Now()

		args := []any{
			"http.client.host",
			req.URL.Host,
			"http.client.path",
			req.URL.Path,
		}

		if body {
			args = append(args, "http.client.request", fmt.Sprintf("%+v", req.Body))
		}
		l.Info(ctx, "http.client: sending request", args...)

		var err error
		defer func() {
			args = append(args, "http.client.latency", time.Since(start).String())
			if err != nil {
				args = append(args, "error", fmt.Sprintf("%+v", err))
			}
			l.Info(ctx, "http.client: received response", args)
		}()

		resp, err := rt.RoundTrip(req)
		if err != nil {
			return nil, err
		}

		args = append(args, "http.client.status", resp.StatusCode)
		if body {
			b, err := httputil.DumpResponse(resp, true)
			if err != nil {
				return resp, fmt.Errorf("dump http response: %+v", err)
			}
			args = append(args, "http.client.response", string(b))
		}
		return resp, nil
	})
}

// todo: implement metricsRoundTripper
func metricsRoundTripper(rt http.RoundTripper) http.RoundTripper {
	return roundTripperFn(func(req *http.Request) (*http.Response, error) {
		return rt.RoundTrip(req)
	})
}
