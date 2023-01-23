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

func resourceSplitSegment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSplitSegmentCreate,
		ReadContext:   resourceSplitSegmentRead,
		DeleteContext: resourceSplitSegmentDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceSplitSegmentImport,
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

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceSplitSegmentImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*Config).API

	importID, parseErr := parseCompositeID(d.Id(), 2)
	if parseErr != nil {
		return nil, parseErr
	}

	workspaceID := importID[0]
	segmentName := importID[1]

	s, _, getErr := client.Segments.Get(workspaceID, segmentName)
	if getErr != nil {
		return nil, getErr
	}

	d.SetId(s.GetName())
	d.Set("workspace_id", workspaceID)
	d.Set("traffic_type_id", s.GetTrafficType().GetID())
	d.Set("name", s.GetName())
	d.Set("description", s.GetDescription())

	return []*schema.ResourceData{d}, nil
}

func resourceSplitSegmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*Config).API
	opts := &api.SegmentRequest{}
	workspaceID := getWorkspaceID(d)

	var trafficTypeID string
	if v, ok := d.GetOk("traffic_type_id"); ok {
		trafficTypeID = v.(string)
		log.Printf("[DEBUG] new segment traffic_type_id is : %v", trafficTypeID)
	}

	if v, ok := d.GetOk("name"); ok {
		opts.Name = v.(string)
		log.Printf("[DEBUG] new segment name is : %v", opts.Name)
	}

	if v, ok := d.GetOk("description"); ok {
		opts.Description = v.(string)
		log.Printf("[DEBUG] new segment description is : %v", opts.Description)
	}

	log.Printf("[DEBUG] Creating segment %s", opts.Name)

	s, _, createErr := client.Segments.Create(workspaceID, trafficTypeID, opts)
	if createErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Unable to create segment %v", opts.Name),
			Detail:   createErr.Error(),
		})
		return diags
	}

	log.Printf("[DEBUG] Created segment %s", opts.Name)

	d.SetId(s.GetName())

	return resourceSplitSegmentRead(ctx, d, meta)
}

func resourceSplitSegmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*Config).API
	workspaceID := getWorkspaceID(d)

	s, _, getErr := client.Segments.Get(workspaceID, d.Id())
	if getErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("unable to fetch segment %s", d.Id()),
			Detail:   getErr.Error(),
		})
		return diags
	}

	d.Set("workspace_id", workspaceID)
	d.Set("traffic_type_id", s.GetTrafficType().GetID())
	d.Set("name", s.GetName())
	d.Set("description", s.GetDescription())

	return diags
}

func resourceSplitSegmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*Config).API
	workspaceID := getWorkspaceID(d)

	log.Printf("[DEBUG] Deleting segment %s", d.Id())

	_, deleteErr := client.Segments.Delete(workspaceID, d.Id())
	if deleteErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("unable to delete segment %s", d.Id()),
			Detail:   deleteErr.Error(),
		})
		return diags
	}

	log.Printf("[DEBUG] Deleted segment %s", d.Id())

	d.SetId("")

	return diags
}
