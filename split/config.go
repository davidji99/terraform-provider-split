package split

import (
	"fmt"
	"github.com/davidji99/terraform-provider-split/api"
	"github.com/davidji99/terraform-provider-split/version"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
)

var (
	UserAgent = fmt.Sprintf("terraform-provider-split/v%s", version.ProviderVersion)
)

type Config struct {
	API     *api.Client
	Headers map[string]string

	apiKey     string
	apiBaseURL string

	RemoveEnvFromStateOnly bool
}

func NewConfig() *Config {
	return &Config{
		RemoveEnvFromStateOnly: false,
	}
}

func (c *Config) initializeAPI() error {
	api, clientInitErr := api.New(api.APIKey(c.apiKey),
		api.APIBaseURL(c.apiBaseURL), api.UserAgent(UserAgent),
		api.CustomHTTPHeaders(c.Headers))
	if clientInitErr != nil {
		return clientInitErr
	}

	c.API = api

	log.Printf("[INFO] Split Client configured")

	return nil
}

func (c *Config) applySchema(d *schema.ResourceData) (err error) {
	if v, ok := d.GetOk("headers"); ok {
		headersRaw := v.(map[string]interface{})
		h := make(map[string]string)

		for k, v := range headersRaw {
			h[k] = fmt.Sprintf("%v", v)
		}

		c.Headers = h
	}

	if v, ok := d.GetOk("base_url"); ok {
		vs := v.(string)
		c.apiBaseURL = vs
	}

	if v, ok := d.GetOk("remove_environment_from_state_only"); ok {
		c.RemoveEnvFromStateOnly = v.(bool)
	}

	return nil
}
