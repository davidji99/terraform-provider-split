package split

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/pmcjury/terraform-provider-split/api"
)

func resourceSplitTrafficTypeAttribute() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSplitTrafficTypeAttributeCreate,
		ReadContext:   resourceSplitTrafficTypeAttributeRead,
		UpdateContext: resourceSplitTrafficTypeAttributeUpdate,
		DeleteContext: resourceSplitTrafficTypeAttributeDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceSplitTrafficTypeAttributeImport,
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

			"identifier": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 200),
			},

			"display_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 200),
			},

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(0, 500),
			},

			"data_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"STRING", "DATETIME", "NUMBER", "SET"}, false),
			},

			"suggested_values": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringLenBetween(1, 50),
				},
				Optional: true,
				Computed: true,
				MaxItems: 50,
			},

			"is_searchable": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"organization_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSplitTrafficTypeAttributeImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*Config).API

	importID, parseErr := parseCompositeID(d.Id(), 3)
	if parseErr != nil {
		return nil, parseErr
	}

	workspaceID := importID[0]
	trafficTypeID := importID[1]
	attributeID := importID[2]

	a, _, getErr := client.Attributes.FindByID(workspaceID, trafficTypeID, attributeID, &api.AttributeListQueryParams{MarkerLimit: 200})
	if getErr != nil {
		return nil, getErr
	}

	d.SetId(a.GetID())
	d.Set("identifier", a.GetID())
	d.Set("workspace_id", workspaceID)
	d.Set("traffic_type_id", trafficTypeID)
	d.Set("display_name", a.GetDisplayName())
	d.Set("description", a.GetDescription())
	d.Set("data_type", a.GetDataType())
	d.Set("is_searchable", a.GetIsSearchable())
	d.Set("organization_id", a.GetOrganizationId())

	return []*schema.ResourceData{d}, nil
}

func resourceSplitTrafficTypeAttributeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*Config).API
	workspaceID := getWorkspaceID(d)
	trafficTypeID := getTrafficTypeID(d)

	opts := constructAttributeRequestOpts(d)
	opts.TrafficTypeID = &trafficTypeID

	log.Printf("[DEBUG] Creating traffic type attribute %v", *opts.Identifier)

	a, _, createErr := client.Attributes.Create(workspaceID, trafficTypeID, opts)
	if createErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Unable to traffic type create attribute %v", opts.DisplayName),
			Detail:   createErr.Error(),
		})
		return diags
	}

	log.Printf("[DEBUG] Created traffic type attribute %v", a.GetID())

	d.SetId(a.GetID())
	d.Set("workspace_id", workspaceID)
	d.Set("traffic_type_id", a.GetTrafficTypeID())

	return resourceSplitTrafficTypeAttributeRead(ctx, d, meta)
}

func resourceSplitTrafficTypeAttributeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*Config).API
	workspaceID := getWorkspaceID(d)
	trafficTypeID := getTrafficTypeID(d)

	a, _, getErr := client.Attributes.FindByID(workspaceID, trafficTypeID, d.Id(), &api.AttributeListQueryParams{MarkerLimit: 200})
	if getErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("unable to fetch traffic type attribute %s", d.Id()),
			Detail:   getErr.Error(),
		})
		return diags
	}

	d.Set("identifier", a.GetID())
	d.Set("display_name", a.GetDisplayName())
	d.Set("description", a.GetDescription())
	d.Set("data_type", a.GetDataType())
	d.Set("is_searchable", a.GetIsSearchable())
	d.Set("organization_id", a.GetOrganizationId())

	if a.HasSuggestedValues() {
		d.Set("suggested_values", a.SuggestedValues)
	}

	return diags
}

func resourceSplitTrafficTypeAttributeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*Config).API
	opts := &api.AttributeRequest{}
	workspaceID := getWorkspaceID(d)
	trafficTypeID := getTrafficTypeID(d)

	if ok := d.HasChange("display_name"); ok {
		vs := d.Get("display_name").(string)
		opts.DisplayName = &vs
		log.Printf("[DEBUG] updated traffic type attribute %v display_name: %v", d.Id(), *opts.DisplayName)
	}

	if ok := d.HasChange("description"); ok {
		_, n := d.GetChange("description")
		var vs string
		if n == nil {
			opts.Description = nil
		} else {
			vs = n.(string)
			opts.Description = &vs
		}

		opts.Description = &vs
		log.Printf("[DEBUG] updated traffic type attribute %v description: %v", d.Id(), *opts.Description)
	}

	if ok := d.HasChange("data_type"); ok {
		_, n := d.GetChange("data_type")
		var vs string
		if n == nil {
			opts.DataType = nil
		} else {
			vs = n.(string)
			opts.DataType = &vs
		}

		opts.DataType = &vs
		log.Printf("[DEBUG] updated traffic type attribute %v data_type: %v", d.Id(), *opts.DataType)
	}

	if ok := d.HasChange("suggested_values"); ok {
		svRaw := d.Get("suggested_values").(*schema.Set).List()
		suggestedValues := make([]string, 0)

		for _, sv := range svRaw {
			suggestedValues = append(suggestedValues, sv.(string))

		}

		opts.SuggestedValues = suggestedValues
		log.Printf("[DEBUG] updated traffic type attribute %v suggested_values: %v", d.Id(), opts.SuggestedValues)
	}

	if ok := d.HasChange("is_searchable"); ok {
		v := d.Get("is_searchable").(bool)
		opts.IsSearchable = &v
		log.Printf("[DEBUG] updated traffic type attribute %v is_searchable: %v", d.Id(), *opts.IsSearchable)
	}

	log.Printf("[DEBUG] Updating traffic type attribute %v", d.Id())

	_, _, updateErr := client.Attributes.Update(workspaceID, trafficTypeID, d.Id(), opts)
	if updateErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Unable to update traffic type attribute %v", d.Id()),
			Detail:   updateErr.Error(),
		})
		return diags
	}

	log.Printf("[DEBUG] Updated traffic type attribute %v", d.Id())

	return resourceSplitTrafficTypeAttributeRead(ctx, d, meta)
}

func resourceSplitTrafficTypeAttributeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Config).API
	var diags diag.Diagnostics
	workspaceID := getWorkspaceID(d)
	trafficTypeID := getTrafficTypeID(d)

	log.Printf("[DEBUG] Deleting attribute %s", d.Id())
	_, deleteErr := client.Attributes.Delete(workspaceID, trafficTypeID, d.Id())
	if deleteErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("unable to delete traffic type attribute %s", d.Id()),
			Detail:   deleteErr.Error(),
		})
		return diags
	}

	log.Printf("[DEBUG] Deleted traffic type attribute %s", d.Id())

	d.SetId("")

	return diags
}

func constructAttributeRequestOpts(d *schema.ResourceData) *api.AttributeRequest {
	opts := &api.AttributeRequest{}

	if v, ok := d.GetOk("identifier"); ok {
		vs := v.(string)
		opts.Identifier = &vs
		log.Printf("[DEBUG] new traffic type attribute identifier is : %v", *opts.Identifier)
	}

	if v, ok := d.GetOk("display_name"); ok {
		vs := v.(string)
		opts.DisplayName = &vs
		log.Printf("[DEBUG] new traffic type attribute display_name is : %v", *opts.DisplayName)
	}

	if v, ok := d.GetOk("description"); ok {
		vs := v.(string)
		opts.Description = &vs
		log.Printf("[DEBUG] new traffic type attribute description is : %v", *opts.Description)
	}

	if v, ok := d.GetOk("data_type"); ok {
		vs := v.(string)
		opts.DataType = &vs
		log.Printf("[DEBUG] new traffic type attribute data_type is : %v", *opts.DataType)
	}

	if v, ok := d.GetOk("suggested_values"); ok {
		svRaw := v.(*schema.Set).List()
		suggestedValues := make([]string, 0)

		for _, sv := range svRaw {
			suggestedValues = append(suggestedValues, sv.(string))

		}
		opts.SuggestedValues = suggestedValues
		log.Printf("[DEBUG] new traffic type attribute suggested_values is : %v", opts.SuggestedValues)
	}

	isSearchable := d.Get("is_searchable").(bool)
	opts.IsSearchable = &isSearchable
	log.Printf("[DEBUG] new traffic type attribute is_searchable is : %v", *opts.IsSearchable)

	return opts
}
