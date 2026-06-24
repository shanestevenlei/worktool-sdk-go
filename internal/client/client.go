// Package client provides a low-level HTTP client for the WorkTool API.
//
// The HTTPClient is intentionally lightweight: it carries only the configuration
// required to issue a single request (base URL + per-request auth) and does NOT
// hold any persistent connection state. Callers (services) are expected to
// instantiate a fresh HTTPClient per request via New().
//
// This design lets the top-level SDK Client stay stateless — multiple goroutines
// sharing the same SDK Client will never contend on a shared HTTPClient, because
// each service method builds its own HTTPClient with its own configuration.
package client

import (
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

// Config carries the connection parameters for a single request.
// All fields are safe to set on every call; nothing is cached.
type Config struct {
	BaseURL  string
	RobotID  string
	HTTPDoer HTTPDoer // optional: if set, replaces resty with this transport (used in tests)
}

// HTTPClient is a one-shot HTTP wrapper around resty.Client.
// Treat instances as request-scoped; do not retain across calls.
type HTTPClient struct {
	cfg Config
	resty *resty.Client
}

// HTTPDoer abstracts resty's underlying http.RoundTripper.
// Use this in tests to substitute the transport without a real network call.
type HTTPDoer interface {
	http.RoundTripper
}

// New constructs a fresh HTTPClient ready for a single request.
func New(cfg Config) *HTTPClient {
	c := resty.New().
		SetBaseURL(cfg.BaseURL).
		SetHeader("Content-Type", "application/json").
		SetTimeout(30 * time.Second).
		SetRetryCount(1)

	if cfg.RobotID != "" {
		c.SetQueryParam("robotId", cfg.RobotID)
	}

	if cfg.HTTPDoer != nil {
		c.SetTransport(cfg.HTTPDoer)
	}

	return &HTTPClient{cfg: cfg, resty: c}
}

// NewWithBaseURL is a convenience helper for the default WorkTool endpoint.
func NewWithBaseURL(robotID string) *HTTPClient {
	return New(Config{
		BaseURL: DefaultBaseURL,
		RobotID: robotID,
	})
}

// DoPOST performs a POST and decodes the response.
func (h *HTTPClient) DoPOST(path string, body, resp interface{}) error {
	_, err := h.resty.R().
		SetBody(body).
		SetResult(resp).
		Post(path)
	return err
}

// DoPOSTRaw performs a POST with a raw []byte body (e.g. JSON arrays).
func (h *HTTPClient) DoPOSTRaw(path string, raw []byte, resp interface{}) error {
	_, err := h.resty.R().
		SetHeader("Content-Type", "application/json").
		SetBody(raw).
		SetResult(resp).
		Post(path)
	return err
}

// DoGET performs a GET with optional query params and decodes the response.
func (h *HTTPClient) DoGET(path string, queryParams map[string]string, resp interface{}) error {
	req := h.resty.R().SetResult(resp)
	for k, v := range queryParams {
		req.SetQueryParam(k, v)
	}
	_, err := req.Get(path)
	return err
}

// Config returns the client configuration (read-only copy).
func (h *HTTPClient) Config() Config {
	return h.cfg
}

// DefaultBaseURL is the production endpoint for the WorkTool API.
const DefaultBaseURL = "https://api.worktool.ymdyes.cn"