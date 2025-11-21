package split

import (
	"context"
	"fmt"
	"github.com/davidji99/terraform-provider-split/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
)

func resourceSplitFlagSet() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSplitFlagSetCreate,
		ReadContext:   resourceSplitFlagSetRead,
		DeleteContext: resourceSplitFlagSetDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceSplitFlagSetImport,
		},

		Schema: map[string]*schema.Schema{
			"workspace_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 50),
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceSplitFlagSetImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*Config).API

	fs, _, getErr := client.FlagSets.FindByID(d.Id())
	if getErr != nil {
		return nil, getErr
	}

	d.SetId(fs.GetID())
	d.Set("workspace_id", fs.Workspace.GetID())
	d.Set("name", fs.GetName())
	d.Set("description", fs.GetDescription())

	return []*schema.ResourceData{d}, nil
}

func resourceSplitFlagSetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*Config).API
	opts := &api.FlagSetRequest{}
	workspaceID := getWorkspaceID(d)

	if v, ok := d.GetOk("name"); ok {
		vs := v.(string)
		opts.Name = &vs
		log.Printf("[DEBUG] new flag set name is : %v", opts.GetName())
	}

	if v, ok := d.GetOk("description"); ok {
		vs := v.(string)
		opts.Description = &vs
		log.Printf("[DEBUG] new flag set description is : %v", opts.GetDescription())
	}

	workspaceType := "WORKSPACE"
	opts.Workspace = &api.WorkspaceIDRef{
		Type: &workspaceType,
		ID:   &workspaceID,
	}
	log.Printf("[DEBUG] new flag set workspace is : %v", workspaceID)

	log.Printf("[DEBUG] Creating Flag Set named %v", opts.GetName())

	fs, _, createErr := client.FlagSets.Create(opts)
	if createErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Unable to create flag set %v", opts.GetName()),
			Detail:   createErr.Error(),
		})
		return diags
	}

	log.Printf("[DEBUG] Created Flag Set named %v", opts.GetName())

	d.SetId(fs.GetID())

	return resourceSplitFlagSetRead(ctx, d, meta)
}

func resourceSplitFlagSetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*Config).API

	fs, _, getErr := client.FlagSets.FindByID(d.Id())
	if getErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("unable to fetch flag set %s", d.Id()),
			Detail:   getErr.Error(),
		})
		return diags
	}

	d.Set("workspace_id", fs.Workspace.GetID())
	d.Set("name", fs.GetName())
	d.Set("description", fs.GetDescription())

	return diags
}

func resourceSplitFlagSetDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*Config).API

	log.Printf("[DEBUG] Deleting Flag Set %s", d.Id())
	_, deleteErr := client.FlagSets.Delete(d.Id())
	if deleteErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("unable to delete flag set %s", d.Id()),
			Detail:   deleteErr.Error(),
		})
		return diags
	}

	log.Printf("[DEBUG] Deleted Flag Set %s", d.Id())

	d.SetId("")

	return diags
}
