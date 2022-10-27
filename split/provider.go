package split

import (
	"context"
	"log"

	"github.com/davidji99/terraform-provider-split/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func New() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SPLIT_API_KEY", nil),
			},

			"base_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SPLIT_API_URL", api.DefaultAPIBaseURL),
			},

			"headers": {
				Type:     schema.TypeMap,
				Elem:     schema.TypeString,
				Optional: true,
			},

			"client_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  300, // 5 min
			},

			"remove_environment_from_state_only": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"split_environment":  dataSourceSplitEnvironment(),
			"split_traffic_type": dataSourceSplitTrafficType(),
			"split_workspace":    dataSourceSplitWorkspace(),
		},

		ResourcesMap: map[string]*schema.Resource{
			//"split_attribute":    resourceSplitAttribute(),
			"split_environment":      resourceSplitEnvironment(),
			"split_group":            resourceSplitGroup(),
			"split_segment":          resourceSplitSegment(),
			"split_split":            resourceSplitSplit(),
			"split_split_definition": resourceSplitSplitDefinition(),
			"split_traffic_type":     resourceSplitTrafficType(),
			"split_user":             resourceSplitUser(),
			"split_workspace":        resourceSplitWorkspace(),
		},

		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	log.Println("[INFO] Initializing Split Provider")

	var diags diag.Diagnostics

	config := NewConfig()

	if applySchemaErr := config.applySchema(d); applySchemaErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to retrieve and set provider attributes",
			Detail:   applySchemaErr.Error(),
		})

		return nil, diags
	}

	if apiKey, ok := d.GetOk("api_key"); ok {
		config.apiKey = apiKey.(string)
	}

	if clientTimeout, ok := d.GetOk("client_timeout"); ok {
		config.clientTimeout = clientTimeout.(int)
	}

	if err := config.initializeAPI(); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to initialize API client",
			Detail:   err.Error(),
		})

		return nil, diags
	}

	log.Printf("[DEBUG] Split Provider initialized")

	return config, diags
}
