package client

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// roundTripFunc adapts a function to http.RoundTripper for tests.
type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestNew(t *testing.T) {
	t.Run("default config", func(t *testing.T) {
		c := New(Config{BaseURL: "https://example.com", RobotID: "robot1"})
		assert.NotNil(t, c)
		assert.Equal(t, "robot1", c.cfg.RobotID)
		assert.Equal(t, "https://example.com", c.cfg.BaseURL)
		assert.NotNil(t, c.resty)
	})

	t.Run("with custom transport", func(t *testing.T) {
		called := false
		transport := roundTripFunc(func(req *http.Request) (*http.Response, error) {
			called = true
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(`{"code":0}`)),
				Header:     http.Header{},
			}, nil
		})
		c := New(Config{BaseURL: "https://example.com", HTTPDoer: transport})

		var resp struct {
			Code int `json:"code"`
		}
		err := c.DoGET("/test", nil, &resp)
		assert.NoError(t, err)
		assert.True(t, called, "transport should have been called")
		assert.Equal(t, 0, resp.Code)
	})
}

func TestNewWithBaseURL(t *testing.T) {
	c := NewWithBaseURL("robot_x")
	assert.Equal(t, DefaultBaseURL, c.cfg.BaseURL)
	assert.Equal(t, "robot_x", c.cfg.RobotID)
}

func TestHTTPClientDoPOST(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var capturedBody []byte
		transport := roundTripFunc(func(req *http.Request) (*http.Response, error) {
			buf, _ := io.ReadAll(req.Body)
			capturedBody = buf
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(`{"code":0,"message":"ok"}`)),
				Header:     http.Header{},
			}, nil
		})
		c := New(Config{HTTPDoer: transport})

		body := map[string]string{"key": "value"}
		var resp struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		}
		err := c.DoPOST("/test", body, &resp)
		assert.NoError(t, err)
		assert.Equal(t, 0, resp.Code)
		assert.Contains(t, string(capturedBody), `"key":"value"`)
	})
}

func TestHTTPClientDoPOSTRaw(t *testing.T) {
	transport := roundTripFunc(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(`{"code":0}`)),
			Header:     http.Header{},
		}, nil
	})
	c := New(Config{HTTPDoer: transport})
	var resp struct {
		Code int `json:"code"`
	}
	err := c.DoPOSTRaw("/test", []byte(`[1,2,3]`), &resp)
	assert.NoError(t, err)
}

func TestHTTPClientDoGET(t *testing.T) {
	t.Run("with query params", func(t *testing.T) {
		var capturedQuery string
		transport := roundTripFunc(func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.RawQuery
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(`{"code":0}`)),
				Header:     http.Header{},
			}, nil
		})
		c := New(Config{HTTPDoer: transport})

		var resp struct {
			Code int `json:"code"`
		}
		err := c.DoGET("/test", map[string]string{"page": "1", "size": "20"}, &resp)
		assert.NoError(t, err)
		assert.Contains(t, capturedQuery, "page=1")
		assert.Contains(t, capturedQuery, "size=20")
	})
}

func TestHTTPClientConfig(t *testing.T) {
	c := New(Config{BaseURL: "https://x", RobotID: "r1"})
	cfg := c.Config()
	assert.Equal(t, "r1", cfg.RobotID)
	assert.Equal(t, "https://x", cfg.BaseURL)
}