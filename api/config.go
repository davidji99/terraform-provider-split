package api

// Config represents all configuration options available to user to customize the API v2.
type Config struct {
	// APIBaseURL is the base URL for Sendgrid's API v3.
	APIBaseURL string

	// UserAgent used when communicating with the Sendgrid API.
	UserAgent string

	// CustomHTTPHeaders are any additional user defined headers.
	CustomHTTPHeaders map[string]string

	// ContentTypeHeader
	ContentTypeHeader string

	// AcceptHeader
	AcceptHeader string

	// APIKey
	APIKey string

	// HarnessToken is used for x-api-key header authentication
	HarnessToken string

	// ClientTimeout
	ClientTimeout int
}

// ParseOptions parses the supplied options functions.
func (c *Config) ParseOptions(opts ...Option) error {
	// Range over each options function and apply it to our API type to
	// configure it. Options functions are applied in order, with any
	// conflicting options overriding earlier calls.
	for _, option := range opts {
		err := option(c)
		if err != nil {
			return err
		}
	}

	return nil
}
