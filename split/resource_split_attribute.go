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

			"attribute_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
				ForceNew: true,
			},
		},
	}
}

func resourceSplitAttributeImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*Config).API

	importID, parseErr := parseCompositeID(d.Id(), 2)
	if parseErr != nil {
		return nil, parseErr
	}

	workspaceID := importID[0]
	envName := importID[1]

	envs, _, getErr := client.Environments.List(workspaceID)
	if getErr != nil {
		return nil, getErr
	}

	notFound := true
	for _, e := range envs {
		if e.GetName() == envName {
			notFound = false
			d.SetId(fmt.Sprintf("%s:%s", workspaceID, e.GetID()))

			d.Set("workspace_id", workspaceID)
			d.Set("name", e.GetName())
			d.Set("production", e.GetProduction())
		}
	}

	if notFound {
		return nil, fmt.Errorf("could not find environment [%s]", envName)
	}

	return []*schema.ResourceData{d}, nil
}

func resourceSplitAttributeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*Config).API
	opts := &api.AttributeRequest{}
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

	d.SetId(fmt.Sprintf("%s:%s", workspaceID, e.GetID()))

	return resourceSplitAttributeRead(ctx, d, meta)
}

func resourceSplitAttributeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*Config).API

	result, parseErr := parseCompositeID(d.Id(), 2)
	if parseErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "unable to parse resource ID during state refresh",
			Detail:   parseErr.Error(),
		})
		return diags
	}

	workspaceID := result[0]
	envID := result[1]

	e, _, getErr := client.Environments.FindByID(workspaceID, envID)
	if getErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("unable to fetch environment %s", envID),
			Detail:   getErr.Error(),
		})
		return diags
	}

	d.Set("workspace_id", workspaceID)
	d.Set("name", e.GetName())
	d.Set("production", e.GetProduction())

	return diags
}

func resourceSplitAttributeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	var diags diag.Diagnostics

	if !config.RemoveEnvFromStateOnly {
		client := config.API

		result, parseErr := parseCompositeID(d.Id(), 2)
		if parseErr != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "unable to parse resource ID during deletion",
				Detail:   parseErr.Error(),
			})
			return diags
		}

		workspaceID := result[0]
		envID := result[1]

		log.Printf("[DEBUG] Deleting Environment %s", envID)
		_, deleteErr := client.Environments.Delete(workspaceID, envID)
		if deleteErr != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("unable to delete environment %s", envID),
				Detail:   deleteErr.Error(),
			})
			return diags
		}

		log.Printf("[DEBUG] Deleted Environment %s", envID)
	}

	d.SetId("")

	return diags
}
