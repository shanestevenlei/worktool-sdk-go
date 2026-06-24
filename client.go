package worktool

import (
	"github.com/shanestevenlei/worktool-sdk-go/internal/client"
	"github.com/shanestevenlei/worktool-sdk-go/service"
)

// Option configures a Client.
type Option func(*config)

type config struct {
	robotID  string
	baseURL  string
	httpDoer client.HTTPDoer
}

// Client composes the various service interfaces. It holds NO HTTP state.
//
// A new HTTP client is created per service method invocation, so concurrent
// calls from multiple goroutines are safe by construction.
type Client struct {
	cfg     config
	Message *service.MessageService
	Robot   *service.RobotService
	History *service.HistoryService
}

// New creates a client with the given options.
//
//	worktool.New(worktool.WithRobotID("your_robot_id"))
func New(opts ...Option) *Client {
	cfg := config{baseURL: client.DefaultBaseURL}
	for _, opt := range opts {
		opt(&cfg)
	}
	c := &Client{
		cfg:     cfg,
		Message: service.NewMessageService(),
		Robot:   service.NewRobotService(),
		History: service.NewHistoryService(),
	}
	// Wire each service to this client's factory. They are stateless beyond
	// the factory pointer, so sharing the Client across goroutines is safe.
	c.Message.SetHTTPFactory(c)
	c.Robot.SetHTTPFactory(c)
	c.History.SetHTTPFactory(c)
	return c
}

// WithRobotID sets the robot identifier (sent as robotId query param).
func WithRobotID(robotID string) Option {
	return func(c *config) { c.robotID = robotID }
}

// WithBaseURL overrides the API endpoint (defaults to production).
func WithBaseURL(baseURL string) Option {
	return func(c *config) { c.baseURL = baseURL }
}

// WithHTTPDoer replaces the default HTTP transport (for tests).
func WithHTTPDoer(doer client.HTTPDoer) Option {
	return func(c *config) { c.httpDoer = doer }
}

// RobotID returns the robot identifier.
func (c *Client) RobotID() string { return c.cfg.robotID }

// HTTPClient builds a fresh HTTP client for the given request.
// Services call this on each invocation to avoid sharing state.
func (c *Client) HTTPClient() *client.HTTPClient {
	return client.New(client.Config{
		BaseURL:  c.cfg.baseURL,
		RobotID:  c.cfg.robotID,
		HTTPDoer: c.cfg.httpDoer,
	})
}
