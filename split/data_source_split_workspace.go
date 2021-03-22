package split

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSplitWorkspace() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSplitWorkspaceRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"requires_title_and_comments": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceSplitWorkspaceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).API

	name := d.Get("name").(string)

	workspace, _, findErr := client.Workspaces.FindByName(name)
	if findErr != nil {
		return diag.FromErr(findErr)
	}

	d.SetId(workspace.GetID())
	d.Set("name", workspace.GetName())
	d.Set("type", workspace.GetType())
	d.Set("requires_title_and_comments", workspace.GetRequiresTitleAndComments())

	return nil
}
