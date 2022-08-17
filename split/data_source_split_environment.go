package split

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceSplitEnvironment() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSplitEnvironmentRead,
		Schema: map[string]*schema.Schema{
			"workspace_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"production": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceSplitEnvironmentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).API

	name := d.Get("name").(string)
	workspaceID := d.Get("workspace_id").(string)

	env, _, findErr := client.Environments.FindByName(workspaceID, name)
	if findErr != nil {
		return diag.FromErr(findErr)
	}

	d.SetId(env.GetID())
	d.Set("name", env.GetName())
	d.Set("production", env.GetProduction())

	return nil
}
