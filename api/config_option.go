package api

import (
	"fmt"
)

// Option is a functional option for configuring the API client.
type Option func(*Config) error

// APIBaseURL allows for a custom API v3 base URL.
func APIBaseURL(url string) Option {
	return func(c *Config) error {
		if err := validateBaseURLOption(url); err != nil {
			return err
		}

		c.APIBaseURL = url
		return nil
	}
}

// UserAgent allows for a custom User Agent.
func UserAgent(userAgent string) Option {
	return func(c *Config) error {
		c.UserAgent = userAgent
		return nil
	}
}

// CustomHTTPHeaders allows for additional HTTPHeaders.
func CustomHTTPHeaders(headers map[string]string) Option {
	return func(c *Config) error {
		c.CustomHTTPHeaders = headers
		return nil
	}
}

// APIKey sets the API key for authentication.
func APIKey(token string) Option {
	return func(c *Config) error {
		c.APIKey = token
		return nil
	}
}

// ClientTimeout sets the client max timeout.
func ClientTimeout(duration int) Option {
	return func(c *Config) error {
		c.ClientTimeout = duration
		return nil
	}
}

// validateBaseURLOption ensures that any custom base URLs do not end with a trailing slash.
func validateBaseURLOption(url string) error {
	// Validate that there is no trailing slashes before setting the custom baseURL
	if url[len(url)-1:] == "/" {
		return fmt.Errorf("custom base URL cannot contain a trailing slash")
	}
	return nil
}

// ContentTypeHeader allows for a custom Content-Type header.
func ContentTypeHeader(s string) Option {
	return func(c *Config) error {
		c.ContentTypeHeader = s
		return nil
	}
}

// AcceptHeader allows for a custom Aceept header.
func AcceptHeader(s string) Option {
	return func(c *Config) error {
		c.AcceptHeader = s
		return nil
	}
}

// HarnessToken sets the Harness token for x-api-key header authentication.
func HarnessToken(token string) Option {
	return func(c *Config) error {
		c.HarnessToken = token
		return nil
	}
}
