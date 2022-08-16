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

func resourceSplitEnvironment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSplitEnvironmentCreate,
		ReadContext:   resourceSplitEnvironmentRead,
		UpdateContext: resourceSplitEnvironmentUpdate,
		DeleteContext: resourceSplitEnvironmentDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceSplitEnvironmentImport,
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
				ValidateFunc: validation.StringLenBetween(1, 15),
			},

			"production": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceSplitEnvironmentImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*Config).API

	importID, parseErr := parseCompositeID(d.Id(), 2)
	if parseErr != nil {
		return nil, parseErr
	}

	workspaceID := importID[0]
	envID := importID[1]

	e, _, getErr := client.Environments.FindByID(workspaceID, envID)
	if getErr != nil {
		return nil, getErr
	}

	d.SetId(e.GetID())

	d.Set("workspace_id", workspaceID)
	d.Set("name", e.GetName())
	d.Set("production", e.GetProduction())

	return []*schema.ResourceData{d}, nil
}

func resourceSplitEnvironmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*Config).API
	opts := &api.EnvironmentRequest{}
	workspaceID := getWorkspaceID(d)

	if v, ok := d.GetOk("name"); ok {
		vs := v.(string)
		opts.Name = &vs
		log.Printf("[DEBUG] new environment name is : %v", opts.GetName())
	}

	production := d.Get("production").(bool)
	opts.Production = &production
	log.Printf("[DEBUG] new environment production is : %v", opts.GetProduction())

	log.Printf("[DEBUG] Creating Environment named %v", opts.GetName())

	e, _, createErr := client.Environments.Create(workspaceID, opts)
	if createErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Unable to create environment %v", opts.GetName()),
			Detail:   createErr.Error(),
		})
		return diags
	}

	log.Printf("[DEBUG] Created Environment named %v", opts.GetName())

	d.SetId(e.GetID())

	return resourceSplitEnvironmentRead(ctx, d, meta)
}

func resourceSplitEnvironmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*Config).API

	e, _, getErr := client.Environments.FindByID(getWorkspaceID(d), d.Id())
	if getErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("unable to fetch environment %s", d.Id()),
			Detail:   getErr.Error(),
		})
		return diags
	}

	d.Set("workspace_id", getWorkspaceID(d))
	d.Set("name", e.GetName())
	d.Set("production", e.GetProduction())

	return diags
}

func resourceSplitEnvironmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*Config).API
	opts := &api.EnvironmentRequest{}

	if ok := d.HasChange("name"); ok {
		vs := d.Get("name").(string)
		opts.Name = &vs
		log.Printf("[DEBUG] updated environment name is : %v", opts.GetName())
	}

	production := d.Get("production").(bool)
	opts.Production = &production
	log.Printf("[DEBUG] updated environment production is : %v", opts.GetProduction())

	log.Printf("[DEBUG] Updating environment")

	_, _, updateErr := client.Environments.Update(getWorkspaceID(d), d.Id(), opts)
	if updateErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("unable to update environment %s", d.Id()),
			Detail:   updateErr.Error(),
		})
		return diags
	}

	log.Printf("[DEBUG] Updated environment")

	return resourceSplitEnvironmentRead(ctx, d, meta)
}

func resourceSplitEnvironmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	var diags diag.Diagnostics

	if !config.RemoveEnvFromStateOnly {
		client := config.API

		log.Printf("[DEBUG] Deleting Environment %s", d.Id())
		_, deleteErr := client.Environments.Delete(getWorkspaceID(d), d.Id())
		if deleteErr != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("unable to delete environment %s", d.Id()),
				Detail:   deleteErr.Error(),
			})
			return diags
		}

		log.Printf("[DEBUG] Deleted Environment %s", d.Id())
	}

	d.SetId("")

	return diags
}
