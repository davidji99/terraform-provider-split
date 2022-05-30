package split

import (
	"context"
	"fmt"
	"github.com/davidji99/terraform-provider-split/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"regexp"
)

func resourceSplitSplit() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSplitSplitCreate,
		ReadContext:   resourceSplitSplitRead,
		UpdateContext: resourceSplitSplitUpdate,
		DeleteContext: resourceSplitSplitDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceSplitSplitImport,
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
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-zA-Z][a-zA-Z-_\d]+$`),
					"Name must start with a letter and can contain hyphens, underscores, letters, and numbers"),
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceSplitSplitImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*Config).API

	importID, parseErr := parseCompositeID(d.Id(), 2)
	if parseErr != nil {
		return nil, parseErr
	}

	workspaceID := importID[0]
	splitID := importID[1]

	s, _, getErr := client.Splits.Get(workspaceID, splitID)
	if getErr != nil {
		return nil, getErr
	}

	d.SetId(s.GetID())
	d.Set("workspace_id", workspaceID)
	d.Set("traffic_type_id", s.GetTrafficType().GetID())
	d.Set("name", s.GetName())
	d.Set("description", s.GetDescription())

	return []*schema.ResourceData{d}, nil
}

func resourceSplitSplitCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*Config).API
	opts := &api.SplitCreateRequest{}
	workspaceID := getWorkspaceID(d)
	trafficTypeID := getTrafficTypeID(d)

	if v, ok := d.GetOk("name"); ok {
		opts.Name = v.(string)
		log.Printf("[DEBUG] new split name is : %v", opts.Name)
	}

	if v, ok := d.GetOk("description"); ok {
		opts.Description = v.(string)
		log.Printf("[DEBUG] new split description is : %v", opts.Description)
	}

	log.Printf("[DEBUG] Creating split %v", opts.Name)

	s, _, createErr := client.Splits.Create(workspaceID, trafficTypeID, opts)
	if createErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Unable to create split %v", opts.Name),
			Detail:   createErr.Error(),
		})
		return diags
	}

	log.Printf("[DEBUG] Created split %v", s.GetID())

	d.SetId(s.GetID())
	d.Set("workspace_id", workspaceID)

	return resourceSplitSplitRead(ctx, d, meta)
}

func resourceSplitSplitUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*Config).API
	workspaceID := getWorkspaceID(d)

	if ok := d.HasChange("description"); ok {
		description := d.Get("description").(string)
		log.Printf("[DEBUG] updated split description is : %v", description)

		log.Printf("[DEBUG] Updating split description %v", d.Id())

		_, _, updateErr := client.Splits.UpdateDescription(workspaceID, d.Get("name").(string), description)
		if updateErr != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Unable to update split description %v", description),
				Detail:   updateErr.Error(),
			})
			return diags
		}

		log.Printf("[DEBUG] Updated split description %v", d.Id())
	}

	return resourceSplitSplitRead(ctx, d, meta)
}

func resourceSplitSplitRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*Config).API

	s, _, getErr := client.Splits.Get(getWorkspaceID(d), d.Id())
	if getErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("unable to fetch split %s", d.Id()),
			Detail:   getErr.Error(),
		})
		return diags
	}

	d.Set("name", s.GetName())
	d.Set("description", s.GetDescription())
	d.Set("traffic_type_id", s.GetTrafficType().GetID())

	return diags
}

func resourceSplitSplitDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Config).API
	var diags diag.Diagnostics

	log.Printf("[DEBUG] Deleting split %s", d.Id())

	_, deleteErr := client.Splits.Delete(getWorkspaceID(d), d.Get("name").(string))
	if deleteErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("unable to delete split %s", d.Id()),
			Detail:   deleteErr.Error(),
		})
		return diags
	}

	log.Printf("[DEBUG] Deleted split %s", d.Id())

	d.SetId("")

	return diags
}
