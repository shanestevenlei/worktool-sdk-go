package service

import "github.com/shanestevenlei/worktool-sdk-go/internal/client"

// HTTPClientFactory builds a fresh HTTP client per request.
// Implementations live in the top-level package (see worktool.Client).
//
// Each service calls factory.HTTPClient() at the start of every method,
// ensuring no HTTP state is shared across goroutines or methods.
type HTTPClientFactory interface {
	HTTPClient() *client.HTTPClient
}

// API path constants shared across services.
const (
	sendRawMessagePath = "/wework/sendRawMessage"
	listRawMessagePath = "/wework/listRawMessage"
)
