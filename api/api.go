package api

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/davidji99/simpleresty"
)

const (
	// DefaultAPIBaseURL is the base API url.
	DefaultAPIBaseURL = "https://api.split.io/internal/api/v2"

	// DefaultUserAgent is the user agent used when making API calls.
	DefaultUserAgent = "split-go"

	// DefaultContentTypeHeader is the default and Content-Type header.
	DefaultContentTypeHeader = "application/json"

	// DefaultAcceptHeader is the default and Content-Type header.
	DefaultAcceptHeader = "application/json"

	//DefaultClientTimeout is the default timeout before the client stops making api calls, will be overriden by split/config.go when instantiate
	DefaultClientTimeout = 300 // 5 min

	UserStatusPending     = "PENDING"
	UserStatusActive      = "ACTIVE"
	UserStatusDeactivated = "DEACTIVATED"

	timeoutError = "reached maximum client timeout"
)

// Client manages communication with Sendgrid APIs.
type Client struct {
	// HTTP client used to communicate with the API.
	http *simpleresty.Client

	// Reuse a single struct instead of allocating one for each service on the heap.
	common service

	// config represents all of the API's configurations.
	config *Config

	//expiresAt is a property that contains the time when the timeout happens
	expiresAt time.Time

	// Services used for talking to different parts of the Sendgrid APIv3.
	ApiKeys      *KeysService
	Attributes   *AttributesService
	Environments *EnvironmentsService
	Groups       *GroupsService
	TrafficTypes *TrafficTypesService
	Segments     *SegmentsService
	Splits       *SplitsService
	Users        *UsersService
	Workspaces   *WorkspacesService
}

// service represents the API service client.
type service struct {
	client *Client
}

// GenericListResult is the generic list result.
type GenericListResult struct {
	Offset     *int `json:"offset"`
	Limit      *int `json:"limit"`
	TotalCount *int `json:"totalCount"`
}

// GenericListQueryParams are parameters for any resource.
type GenericListQueryParams struct {
	// The offset to retrieve. Useful for pagination
	Offset int `url:"offset,omitempty"`

	// The maximum limit to return per call. Max=20-50.
	Limit int `url:"limit,omitempty"`
}

// New constructs a new Client.
func New(opts ...Option) (*Client, error) {
	config := &Config{
		APIBaseURL:        DefaultAPIBaseURL,
		UserAgent:         DefaultUserAgent,
		ContentTypeHeader: DefaultContentTypeHeader,
		AcceptHeader:      DefaultAcceptHeader,
		ClientTimeout:     DefaultClientTimeout,
		APIKey:            "",
	}

	// Define any user custom Client settings
	if optErr := config.ParseOptions(opts...); optErr != nil {
		return nil, optErr
	}

	expiresAt := time.Now().Add(time.Duration(config.ClientTimeout) * time.Second)

	client := &Client{
		config:    config,
		http:      simpleresty.NewWithBaseURL(config.APIBaseURL),
		expiresAt: expiresAt,
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
	c.ApiKeys = (*KeysService)(&c.common)
	c.Attributes = (*AttributesService)(&c.common)
	c.Environments = (*EnvironmentsService)(&c.common)
	c.Groups = (*GroupsService)(&c.common)
	c.TrafficTypes = (*TrafficTypesService)(&c.common)
	c.Segments = (*SegmentsService)(&c.common)
	c.Splits = (*SplitsService)(&c.common)
	c.Users = (*UsersService)(&c.common)
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

func (c *Client) checkRateLimit(resp *simpleresty.Response) bool {
	if resp != nil && resp.StatusCode == 429 {
		remainingOrgSeconds, _ := strconv.Atoi((resp.Resp.Header().Get("X-RateLimit-Reset-Seconds-Org")))
		timeToSleep, _ := strconv.Atoi(resp.Resp.Header().Get("X-RateLimit-Reset-Seconds-IP"))
		if remainingOrgSeconds != 0 {
			// Got rate-limit by Organization
			timeToSleep = remainingOrgSeconds
		}
		// Got rate-limit by IP-addr
		log.Printf("[DEBUG] Got rate-limited, sleeping for %d seconds", timeToSleep)
		time.Sleep(time.Duration(timeToSleep) * time.Second)
		return true
	}
	return false
}

// checkTimeout returns true if timeout, false if still have time
func (c *Client) checkTimeout() bool {
	return time.Now().After(c.expiresAt)
}

func (c *Client) get(url string, r, body interface{}) (*simpleresty.Response, error) {
	if !c.checkTimeout() {
		response, getErr := c.http.Get(url, &r, body)
		if c.checkRateLimit(response) {
			response, getErr = c.get(url, &r, body)
		}
		return response, getErr
	}
	return nil, errors.New(timeoutError)
}

func (c *Client) post(url string, r, opts interface{}) (*simpleresty.Response, error) {
	if !c.checkTimeout() {
		response, getErr := c.http.Post(url, &r, opts)
		if c.checkRateLimit(response) {
			response, getErr = c.post(url, &r, opts)
		}
		return response, getErr
	}
	return nil, errors.New(timeoutError)
}

func (c *Client) put(url string, r, opts interface{}) (*simpleresty.Response, error) {
	if !c.checkTimeout() {
		response, getErr := c.http.Put(url, &r, opts)
		if c.checkRateLimit(response) {
			response, getErr = c.put(url, &r, opts)
		}
		return response, getErr
	}
	return nil, errors.New(timeoutError)
}

func (c *Client) patch(url string, r, opts interface{}) (*simpleresty.Response, error) {
	if !c.checkTimeout() {
		response, getErr := c.http.Patch(url, &r, opts)
		if c.checkRateLimit(response) {
			response, getErr = c.patch(url, &r, opts)
		}
		return response, getErr
	}
	return nil, errors.New(timeoutError)
}

func (c *Client) delete(url string, r, opts interface{}) (*simpleresty.Response, error) {
	if !c.checkTimeout() {
		response, getErr := c.http.Delete(url, &r, opts)
		if c.checkRateLimit(response) {
			response, getErr = c.delete(url, &r, opts)
		}
		return response, getErr
	}
	return nil, errors.New(timeoutError)
}
