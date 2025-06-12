package split

import (
	"context"
	"fmt"
	"github.com/davidji99/terraform-provider-split/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
)

func resourceSplitGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSplitGroupCreate,
		ReadContext:   resourceSplitGroupRead,
		UpdateContext: resourceSplitGroupUpdate,
		DeleteContext: resourceSplitGroupDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceSplitGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*Config).API
	opts := &api.GroupRequest{}

	if v, ok := d.GetOk("name"); ok {
		opts.Name = v.(string)
		log.Printf("[DEBUG] new group name is : %v", opts.Name)
	}

	if v, ok := d.GetOk("description"); ok {
		opts.Description = v.(string)
		log.Printf("[DEBUG] new group description is : %v", opts.Description)
	}

	log.Printf("[DEBUG] Creating group %s", opts.Name)

	g, _, createErr := client.Groups.Create(opts)
	if createErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Unable to create group %v", opts.Name),
			Detail:   createErr.Error(),
		})
		return diags
	}

	log.Printf("[DEBUG] Created group %s", opts.Name)

	d.SetId(g.GetID())

	return resourceSplitGroupRead(ctx, d, meta)
}

func resourceSplitGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*Config).API

	g, _, getErr := client.Groups.Get(d.Id())
	if getErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("unable to fetch group %s", d.Id()),
			Detail:   getErr.Error(),
		})
		return diags
	}

	d.Set("name", g.GetName())
	d.Set("description", g.GetDescription())

	return diags
}

func resourceSplitGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*Config).API
	opts := &api.GroupRequest{}

	if ok := d.HasChange("name"); ok {
		opts.Name = d.Get("name").(string)
		log.Printf("[DEBUG] updated group name is : %v", opts.Name)
	}

	if ok := d.HasChange("description"); ok {
		opts.Description = d.Get("description").(string)
		log.Printf("[DEBUG] updated group description is : %v", opts.Description)
	}

	log.Printf("[DEBUG] Updating group %s", d.Id())

	_, _, updateErr := client.Groups.Update(d.Id(), opts)
	if updateErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Unable to update group %v", d.Id()),
			Detail:   updateErr.Error(),
		})
		return diags
	}

	log.Printf("[DEBUG] Updated group %s", d.Id())

	return resourceSplitGroupRead(ctx, d, meta)
}

func resourceSplitGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*Config).API

	log.Printf("[DEBUG] Deleting group %s", d.Id())

	_, deleteErr := client.Groups.Delete(d.Id())
	if deleteErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("unable to delete group %s", d.Id()),
			Detail:   deleteErr.Error(),
		})
		return diags
	}

	log.Printf("[DEBUG] Deleted group %s", d.Id())

	d.SetId("")

	return diags
}

// resourceSplitGroupWithDeprecation wraps resourceSplitGroup and adds deprecation checks for harness_token
func resourceSplitGroupWithDeprecation() *schema.Resource {
	r := resourceSplitGroup()

	// Add plan-time validation using CustomizeDiff
	r.CustomizeDiff = func(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) error {
		if diags := checkResourceDeprecationWithHarnessToken("split_group", meta); len(diags) > 0 {
			return fmt.Errorf(diags[0].Summary + ": " + diags[0].Detail)
		}
		return nil
	}

	// Wrap create function with deprecation check
	originalCreate := r.CreateContext
	r.CreateContext = func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		// Check if the harness_token is set
		if diags := checkResourceDeprecationWithHarnessToken("split_group", meta); len(diags) > 0 {
			return diags
		}

		return originalCreate(ctx, d, meta)
	}

	// Wrap read function with deprecation check
	originalRead := r.ReadContext
	r.ReadContext = func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		// Check if the harness_token is set
		if diags := checkResourceDeprecationWithHarnessToken("split_group", meta); len(diags) > 0 {
			return diags
		}

		return originalRead(ctx, d, meta)
	}

	// Wrap update function with deprecation check
	originalUpdate := r.UpdateContext
	r.UpdateContext = func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		// Check if the harness_token is set
		if diags := checkResourceDeprecationWithHarnessToken("split_group", meta); len(diags) > 0 {
			return diags
		}

		return originalUpdate(ctx, d, meta)
	}

	return r
}
