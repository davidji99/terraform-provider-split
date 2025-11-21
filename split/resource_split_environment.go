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

			"api_token_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "IDs of automatically-created API tokens for this environment",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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

	// Store the automatically-created API token IDs
	if e.ApiTokens != nil && len(e.ApiTokens) > 0 {
		tokenIDs := make([]string, 0, len(e.ApiTokens))
		for _, token := range e.ApiTokens {
			if token.ID != nil {
				tokenIDs = append(tokenIDs, *token.ID)
				log.Printf("[DEBUG] Storing API token ID %s for environment %s", *token.ID, e.GetID())
			}
		}
		d.Set("api_token_ids", tokenIDs)
	}

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

	// Preserve API token IDs in state if they exist
	// The API List endpoint doesn't always return apiTokens, so we keep what's in state
	if e.ApiTokens != nil && len(e.ApiTokens) > 0 {
		tokenIDs := make([]string, 0, len(e.ApiTokens))
		for _, token := range e.ApiTokens {
			if token.ID != nil {
				tokenIDs = append(tokenIDs, *token.ID)
			}
		}
		d.Set("api_token_ids", tokenIDs)
	}
	// If no tokens in response, keep existing state (don't overwrite)

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

		// Delete associated API tokens first (required before environment deletion)
		if tokenIDs, ok := d.GetOk("api_token_ids"); ok {
			tokenList := tokenIDs.([]interface{})
			for _, tokenID := range tokenList {
				tokenIDStr := tokenID.(string)
				log.Printf("[DEBUG] Deleting API token %s for environment %s", tokenIDStr, d.Id())
				_, deleteTokenErr := client.ApiKeys.Delete(tokenIDStr)
				if deleteTokenErr != nil {
					log.Printf("[WARN] Error deleting API token %s: %s", tokenIDStr, deleteTokenErr.Error())
					// Continue trying to delete other tokens and the environment
				} else {
					log.Printf("[DEBUG] Deleted API token %s", tokenIDStr)
				}
			}
		}

		log.Printf("[DEBUG] Deleting Environment %s", d.Id())
		_, deleteErr := client.Environments.Delete(getWorkspaceID(d), d.Id())
		if deleteErr != nil {
			// Provide helpful error message for token-related deletion issues
			errorDetail := deleteErr.Error()
			if _, hasTokens := d.GetOk("api_token_ids"); !hasTokens {
				errorDetail = fmt.Sprintf("%s\n\nNote: Environments cannot be deleted until all associated API tokens are revoked. If this environment was created with an older provider version, you may need to manually delete the auto-created API tokens via the Split.io UI or API before destroying this environment.", errorDetail)
			}
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("unable to delete environment %s", d.Id()),
				Detail:   errorDetail,
			})
			return diags
		}

		log.Printf("[DEBUG] Deleted Environment %s", d.Id())
	}

	d.SetId("")

	return diags
}
