package split

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceSplitFlagSet() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSplitFlagSetRead,
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

			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceSplitFlagSetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).API

	name := d.Get("name").(string)
	workspaceID := d.Get("workspace_id").(string)

	fs, _, findErr := client.FlagSets.FindByName(workspaceID, name)
	if findErr != nil {
		return diag.FromErr(findErr)
	}

	d.SetId(fs.GetID())
	d.Set("name", fs.GetName())
	d.Set("description", fs.GetDescription())

	return nil
}
