package api

import (
	"fmt"
	"github.com/davidji99/simpleresty"
	"time"
)

const (
	// DefaultAPIBaseURL
	DefaultAPIBaseURL = "https://api.split.io/internal/api/v2"

	// DefaultUserAgent is the user agent used when making API calls.
	DefaultUserAgent = "split-go"

	// DefaultContentTypeHeader is the default and Content-Type header.
	DefaultContentTypeHeader = "application/json"

	// DefaultAcceptHeader is the default and Content-Type header.
	DefaultAcceptHeader = "application/json"
)

// Client manages communication with Sendgrid APIs.
type Client struct {
	// HTTP client used to communicate with the API.
	http *simpleresty.Client

	// Reuse a single struct instead of allocating one for each service on the heap.
	common service

	// config represents all of the API's configurations.
	config *Config

	// Services used for talking to different parts of the Sendgrid APIv3.
	Environments *EnvironmentsService
	TrafficTypes *TrafficTypesService
	Segments     *SegmentsService
	Workspaces   *WorkspacesService
}

// service represents the API service client.
type service struct {
	client *Client
}

// GenericListResult
type GenericListResult struct {
	Offset     *int `json:"offset"`
	Limit      *int `json:"limit"`
	TotalCount *int `json:"totalCount"`
}

// New constructs a new Client.
func New(opts ...Option) (*Client, error) {
	config := &Config{
		APIBaseURL:        DefaultAPIBaseURL,
		UserAgent:         DefaultUserAgent,
		ContentTypeHeader: DefaultContentTypeHeader,
		AcceptHeader:      DefaultAcceptHeader,
		APIKey:            "",
	}

	// Define any user custom Client settings
	if optErr := config.ParseOptions(opts...); optErr != nil {
		return nil, optErr
	}

	client := &Client{
		config: config,
		http:   simpleresty.NewWithBaseURL(config.APIBaseURL),
	}

	// Set headers
	client.setHeaders()

	// Inject services
	client.injectServices()

	return client, nil
}

// injectServices adds the services to the client.
func (c *Client) injectServices() {
	c.common.client = c
	c.Environments = (*EnvironmentsService)(&c.common)
	c.TrafficTypes = (*TrafficTypesService)(&c.common)
	c.Segments = (*SegmentsService)(&c.common)
	c.Workspaces = (*WorkspacesService)(&c.common)
}

func (c *Client) setHeaders() {
	c.http.SetHeader("Content-type", c.config.ContentTypeHeader).
		SetHeader("Accept", c.config.AcceptHeader).
		SetHeader("User-Agent", c.config.UserAgent).
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", c.config.APIKey)).
		SetTimeout(2 * time.Minute).
		SetAllowGetMethodPayload(true)

	// Set additional headers
	if c.config.CustomHTTPHeaders != nil {
		c.http.SetHeaders(c.config.CustomHTTPHeaders)
	}
}
