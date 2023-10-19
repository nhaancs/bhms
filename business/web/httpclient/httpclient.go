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
	logBody bool
	proxy   func(request *http.Request) (*url.URL, error)
}

type Option func(o *options)

type roundTripperFn func(req *http.Request) (*http.Response, error)

func (f roundTripperFn) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func WithLogBody() Option {
	return func(o *options) {
		o.logBody = true
	}
}

func WithProxy(proxy *url.URL) Option {
	return func(o *options) {
		if proxy != nil {
			o.proxy = http.ProxyURL(proxy)
		}
	}
}

var transport = &http.Transport{
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
// To log request, response body use WithLogBody
// To add proxy use WithProxy
func New(log *logger.Logger, opts ...Option) *http.Client {
	o := options{
		logBody: false,
		proxy:   nil,
	}
	for _, opt := range opts {
		opt(&o)
	}

	t := transport
	if o.proxy != nil {
		t = t.Clone()
		t.Proxy = o.proxy
	}

	return &http.Client{
		Timeout:   30 * time.Second,
		Transport: metricsRoundTripper(logRoundTripper(otelhttp.NewTransport(t), log, o.logBody)),
	}
}

func logRoundTripper(rt http.RoundTripper, log *logger.Logger, body bool) http.RoundTripper {
	return roundTripperFn(func(req *http.Request) (*http.Response, error) {
		//start := time.Now()

		//l := log.WithFields(map[string]interface{}{
		//	"http.client.host": req.URL.Host,
		//	"http.client.path": req.URL.Path,
		//})

		if body {
			//l.WithField("http.client.request", fmt.Sprintf("%+v", req.Body)).Info("http.client: sending request")
		} else {
			//l.Info("http.client: sending request")
		}

		var err error
		defer func() {
			if err != nil {
				//l = l.WithField("error", err)
			}
			//l.WithField("http.client.latency", time.Since(start).String()).Info("http.client: received response")
		}()

		resp, err := rt.RoundTrip(req)
		if err != nil {
			return nil, err
		}

		//l = l.WithField("http.client.status", resp.StatusCode)
		if !body {
			return resp, nil
		}
		_, err = httputil.DumpResponse(resp, true)
		if err != nil {
			return resp, fmt.Errorf("dump http response: %+v", err)
		}

		//l = l.WithField("http.client.response", string(b))

		return resp, err
	})
}

// todo: implement metricsRoundTripper
func metricsRoundTripper(rt http.RoundTripper) http.RoundTripper {
	return roundTripperFn(func(req *http.Request) (*http.Response, error) {
		return rt.RoundTrip(req)
	})
}
