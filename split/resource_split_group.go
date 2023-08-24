package split

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pmcjury/terraform-provider-split/api"
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
