package httpclient

import (
	"context"
	"fmt"
	"io"
	"maps"
	"net/http"

	"github.com/gojek/heimdall/v7/httpclient"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type HeimdallConfig struct {
	Headers     http.Header
	OtelEnabled bool
}

type Heimdall struct {
	client  *httpclient.Client
	headers http.Header
}

func NewHeimdall(config HeimdallConfig) *Heimdall {

	httpClient := &http.Client{
		Transport: http.DefaultTransport,
	}

	if config.OtelEnabled {
		httpClient.Transport = otelhttp.NewTransport(
			http.DefaultTransport,
			otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
				return fmt.Sprintf("%s %s", r.Method, r.URL.Path)
			}),
		)
	}

	client := httpclient.NewClient(
		httpclient.WithHTTPClient(httpClient),
	)

	return &Heimdall{
		client:  client,
		headers: config.Headers,
	}
}

func (h *Heimdall) Do(req *http.Request) (*http.Response, error) {
	maps.Copy(req.Header, h.headers)
	return h.client.Do(req)
}

func (h *Heimdall) Get(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return h.Do(req)
}

func (h *Heimdall) Post(ctx context.Context, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	return h.Do(req)
}

func (h *Heimdall) Put(ctx context.Context, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, body)
	if err != nil {
		return nil, err
	}
	return h.client.Do(req)
}

func (h *Heimdall) Delete(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	return h.client.Do(req)
}
