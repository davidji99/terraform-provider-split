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

func resourceSplitSegmentEnvironmentAssociation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSplitSegmentEnvironmentAssociationCreate,
		ReadContext:   resourceSplitSegmentEnvironmentAssociationRead,
		DeleteContext: resourceSplitSegmentEnvironmentAssociationDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceSplitSegmentEnvironmentAssociationImport,
		},

		Schema: map[string]*schema.Schema{
			"workspace_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"environment_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"segment_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceSplitSegmentEnvironmentAssociationImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*Config).API

	importID, parseErr := parseCompositeID(d.Id(), 3)
	if parseErr != nil {
		return nil, parseErr
	}

	workspaceID := importID[0]
	environmentID := importID[1]
	segmentName := importID[2]

	segments, _, getErr := client.Environments.ListSegments(workspaceID, environmentID)
	if getErr != nil {
		return nil, fmt.Errorf(fmt.Sprintf("unable to fetch all segments in environment %s", environmentID))
	}

	// Iterate through all segments to find the right one
	var segment *api.Segment
	for _, s := range segments.Objects {
		if s.GetName() == segmentName {
			segment = s
		}
	}

	if segment == nil {
		return nil, fmt.Errorf(fmt.Sprintf("did not find to segment [%s] in environment [%s]", d.Id(), environmentID))
	}

	d.SetId(segment.GetName())
	d.Set("workspace_id", workspaceID)
	d.Set("environment_id", segment.GetEnvironment().GetID())

	return []*schema.ResourceData{d}, nil
}

func resourceSplitSegmentEnvironmentAssociationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*Config).API
	environmentID := getEnvironmentID(d)
	segmentName := d.Get("segment_name").(string)

	log.Printf("[DEBUG] Activating segment [%s] in environment [%s]", segmentName, environmentID)

	s, _, createErr := client.Segments.Activate(environmentID, segmentName)
	if createErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Unable to activate segment [%s] in environment [%s]", segmentName, environmentID),
			Detail:   createErr.Error(),
		})
		return diags
	}

	log.Printf("[DEBUG] Activated segment [%s] in environment [%s]", segmentName, environmentID)

	d.SetId(s.GetName())

	return resourceSplitSegmentEnvironmentAssociationRead(ctx, d, meta)
}

func resourceSplitSegmentEnvironmentAssociationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*Config).API
	workspaceID := getWorkspaceID(d)
	environmentID := getEnvironmentID(d)

	segments, _, getErr := client.Environments.ListSegments(workspaceID, environmentID)
	if getErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("unable to fetch all segments in environment %s", environmentID),
			Detail:   getErr.Error(),
		})
		return diags
	}

	// Iterate through all segments to find the right one
	var segment *api.Segment
	for _, s := range segments.Objects {
		if s.GetName() == d.Id() {
			segment = s
		}
	}

	if segment == nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("did not find to segment [%s] in environment [%s]", d.Id(), environmentID),
			Detail:   getErr.Error(),
		})
		return diags
	}

	d.Set("workspace_id", workspaceID)
	d.Set("environment_id", segment.GetEnvironment().GetID())
	d.Set("segment_name", segment.GetName())

	return diags
}

func resourceSplitSegmentEnvironmentAssociationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*Config).API
	environmentID := getEnvironmentID(d)

	log.Printf("[DEBUG] Deactivating segment [%s] from environment [%s]", d.Id(), environmentID)

	_, deleteErr := client.Segments.Deactivate(environmentID, d.Id())
	if deleteErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("unable to delete segment %s", d.Id()),
			Detail:   deleteErr.Error(),
		})
		return diags
	}

	log.Printf("[DEBUG] Deactivated segment [%s] from environment [%s]", d.Id(), environmentID)

	d.SetId("")

	return diags
}
