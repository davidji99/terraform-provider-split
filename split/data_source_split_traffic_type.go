package split

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceSplitTrafficType() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSplitTrafficTypeRead,
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

			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"display_attribute_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"workspace_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceSplitTrafficTypeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).API

	name := d.Get("name").(string)
	workspaceID := d.Get("workspace_id").(string)

	trafficType, _, findErr := client.TrafficTypes.FindByName(workspaceID, name)
	if findErr != nil {
		return diag.FromErr(findErr)
	}

	d.SetId(trafficType.GetID())
	d.Set("workspace_id", trafficType.GetWorkspace().GetID())
	d.Set("name", trafficType.GetName())
	d.Set("type", trafficType.GetType())
	d.Set("display_attribute_id", trafficType.GetDisplayAttributeID())
	d.Set("workspace_name", trafficType.GetWorkspace().GetName())

	return nil
}
