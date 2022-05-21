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

func resourceSplitTrafficType() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSplitTrafficTypeCreate,
		ReadContext:   resourceSplitTrafficTypeRead,
		DeleteContext: resourceSplitTrafficTypeDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceSplitTrafficTypeImport,
		},

		Schema: map[string]*schema.Schema{
			"workspace_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"display_attribute_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSplitTrafficTypeImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*Config).API

	importID, parseErr := parseCompositeID(d.Id(), 2)
	if parseErr != nil {
		return nil, parseErr
	}

	workspaceID := importID[0]
	trafficTypeID := importID[1]

	tt, _, getErr := client.TrafficTypes.FindByID(workspaceID, trafficTypeID)
	if getErr != nil {
		return nil, getErr
	}

	d.SetId(tt.GetID())
	d.Set("workspace_id", tt.Workspace.GetID())
	d.Set("name", tt.GetName())
	d.Set("type", tt.GetType())
	d.Set("display_attribute_id", tt.GetDisplayAttributeID())

	return []*schema.ResourceData{d}, nil
}

func resourceSplitTrafficTypeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*Config).API
	opts := &api.TrafficTypeRequest{}
	workspaceID := getWorkspaceID(d)

	if v, ok := d.GetOk("name"); ok {
		opts.Name = v.(string)
		log.Printf("[DEBUG] new traffuc type is : %v", opts.Name)
	}

	log.Printf("[DEBUG] Creating traffic type %v", opts.Name)

	tt, _, createErr := client.TrafficTypes.Create(workspaceID, opts)
	if createErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Unable to create traffic type %v", opts.Name),
			Detail:   createErr.Error(),
		})
		return diags
	}

	log.Printf("[DEBUG] Created traffuc type %v", opts.Name)

	d.SetId(tt.GetID())
	d.Set("workspace_id", tt.Workspace.GetID())

	return resourceSplitTrafficTypeRead(ctx, d, meta)
}

func resourceSplitTrafficTypeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*Config).API

	tt, _, getErr := client.TrafficTypes.FindByID(getWorkspaceID(d), d.Id())
	if getErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("unable to fetch traffic type %s", d.Id()),
			Detail:   getErr.Error(),
		})
		return diags
	}

	d.Set("name", tt.GetName())
	d.Set("type", tt.GetType())
	d.Set("display_attribute_id", tt.GetDisplayAttributeID())

	return diags
}

func resourceSplitTrafficTypeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Config).API
	var diags diag.Diagnostics

	log.Printf("[DEBUG] Deleting traffic type %s", d.Id())
	_, deleteErr := client.TrafficTypes.Delete(d.Id())
	if deleteErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("unable to delete traffic type %s", d.Id()),
			Detail:   deleteErr.Error(),
		})
		return diags
	}

	log.Printf("[DEBUG] Deleted traffic type %s", d.Id())

	d.SetId("")

	return diags
}
