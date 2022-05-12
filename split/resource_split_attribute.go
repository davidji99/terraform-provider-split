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

func resourceSplitAttribute() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSplitAttributeCreate,
		ReadContext:   resourceSplitAttributeRead,
		DeleteContext: resourceSplitAttributeDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceSplitAttributeImport,
		},

		Schema: map[string]*schema.Schema{
			"workspace_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"traffic_type_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"display_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"description": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"data_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"string", "datetime", "number", "set"}, false),
			},

			"is_searchable": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
			},
		},
	}
}

func resourceSplitAttributeImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*Config).API

	importID, parseErr := parseCompositeID(d.Id(), 3)
	if parseErr != nil {
		return nil, parseErr
	}

	workspaceID := importID[0]
	trafficTypeID := importID[1]
	attributeID := importID[2]

	a, _, getErr := client.Attributes.FindByID(workspaceID, trafficTypeID, attributeID)
	if getErr != nil {
		return nil, getErr
	}

	d.SetId(a.GetID())
	d.Set("workspace_id", workspaceID)
	d.Set("traffic_type_id", trafficTypeID)
	d.Set("display_name", a.GetDisplayName())
	d.Set("description", a.GetDescription())
	d.Set("data_type", a.GetDataType())
	d.Set("is_searchable", a.GetIsSearchable())

	return []*schema.ResourceData{d}, nil
}

func resourceSplitAttributeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*Config).API
	opts := &api.AttributeRequest{}
	workspaceID := getWorkspaceID(d)
	trafficTypeID := getTrafficTypeID(d)

	if v, ok := d.GetOk("display_name"); ok {
		opts.DisplayName = v.(string)
		log.Printf("[DEBUG] new attribute display_name is : %v", opts.DisplayName)
	}

	if v, ok := d.GetOk("description"); ok {
		opts.Description = v.(string)
		log.Printf("[DEBUG] new attribute description is : %v", opts.Description)
	}

	if v, ok := d.GetOk("data_type"); ok {
		opts.DataType = v.(string)
		log.Printf("[DEBUG] new attribute DataType is : %v", opts.DataType)
	}

	isSearchable := d.Get("is_searchable").(bool)
	opts.IsSearchable = &isSearchable

	log.Printf("[DEBUG] Creating attribute %v", opts.DisplayName)

	e, _, createErr := client.Attributes.Create(workspaceID, trafficTypeID, opts)
	if createErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Unable to create attribute %v", opts.DisplayName),
			Detail:   createErr.Error(),
		})
		return diags
	}

	log.Printf("[DEBUG] Created attribute %v", opts.DisplayName)

	d.SetId(e.GetID())
	d.Set("workspace_id", workspaceID)
	d.Set("traffic_type_id", trafficTypeID)

	return resourceSplitAttributeRead(ctx, d, meta)
}

func resourceSplitAttributeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*Config).API
	workspaceID := getWorkspaceID(d)
	trafficTypeID := getTrafficTypeID(d)

	a, _, getErr := client.Attributes.FindByID(workspaceID, trafficTypeID, d.Id())
	if getErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("unable to fetch attribute %s", d.Id()),
			Detail:   getErr.Error(),
		})
		return diags
	}

	d.Set("display_name", a.GetDisplayName())
	d.Set("description", a.GetDescription())
	d.Set("data_type", a.GetDataType())
	d.Set("is_searchable", a.GetIsSearchable())

	return diags
}

func resourceSplitAttributeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Config).API
	var diags diag.Diagnostics
	workspaceID := getWorkspaceID(d)
	trafficTypeID := getTrafficTypeID(d)

	log.Printf("[DEBUG] Deleting attribute %s", d.Id())
	_, deleteErr := client.Attributes.Delete(workspaceID, trafficTypeID, d.Id())
	if deleteErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("unable to delete attribute %s", d.Id()),
			Detail:   deleteErr.Error(),
		})
		return diags
	}

	log.Printf("[DEBUG] Deleted attribute %s", d.Id())

	d.SetId("")

	return diags
}
