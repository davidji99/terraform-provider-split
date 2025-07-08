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

func resourceSplitApiKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSplitApiKeyCreate,
		ReadContext:   resourceSplitApiKeyRead,
		DeleteContext: resourceSplitApiKeyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceSplitApiKeyImport,
		},

		Schema: map[string]*schema.Schema{
			"workspace_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"environment_ids": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsUUID,
				},
			},

			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 100),
			},

			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(api.ValidApiKeyTypes, false),
			},

			"roles": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(api.ValidApiKeyRoles, false),
				},
			},
		},
	}
}

// Since there's no `GET https://api.split.io/internal/api/v2/apiKeys` endpoint, it is not possible to import existing keys.
func resourceSplitApiKeyImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return nil, fmt.Errorf("not possible to import existing API keys due to API limitations")
}

func resourceSplitApiKeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*Config).API
	opts := &api.KeyRequest{}

	if v, ok := d.GetOk("workspace_id"); ok {
		vs := v.(string)
		log.Printf("[DEBUG] new api key workspace_id: %s", vs)
		opts.Workspace = &api.KeyWorkspaceRequest{
			Type: "workspace",
			Id:   vs,
		}
	}

	if v, ok := d.GetOk("name"); ok {
		vs := v.(string)
		opts.Name = vs
		log.Printf("[DEBUG] new api key name is : %v", vs)
	}

	if v, ok := d.GetOk("type"); ok {
		vs := v.(string)
		opts.KeyType = vs
		log.Printf("[DEBUG] new api key type is : %v", vs)
	}

	if v, ok := d.GetOk("roles"); ok {
		rolesListRaw := v.([]interface{})
		roleList := make([]string, 0)
		for _, roleRaw := range rolesListRaw {
			roleList = append(roleList, roleRaw.(string))
		}
		opts.Roles = roleList
		log.Printf("[DEBUG] new api key roles is : %v", roleList)
	}

	if v, ok := d.GetOk("environment_ids"); ok {
		envIdsListRaw := v.([]interface{})
		envIdsMapList := make([]api.KeyEnvironmentRequest, 0)
		for _, envIdRaw := range envIdsListRaw {
			envIdsMapList = append(envIdsMapList, api.KeyEnvironmentRequest{
				Type: "environment",
				Id:   envIdRaw.(string),
			})
		}
		opts.Environments = envIdsMapList
		log.Printf("[DEBUG] new api key environments is : %v", envIdsMapList)
	}

	log.Printf("[DEBUG] Creating new api key %v", opts.Name)

	apiKey, _, createErr := client.ApiKeys.Create(opts)
	if createErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Unable to create new api key %v", opts.Name),
			Detail:   createErr.Error(),
		})
		return diags
	}

	log.Printf("[DEBUG] Created new api key %v", apiKey.GetName())

	d.SetId(apiKey.GetKey())

	return resourceSplitApiKeyRead(ctx, d, meta)
}

// Since there's no `GET https://api.split.io/internal/api/v2/apiKeys` endpoint, resourceSplitApiKeyRead will be a `noop`.
func resourceSplitApiKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[WARN] can't refresh API key resource due to API limiitations")
	return nil
}

func resourceSplitApiKeyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Config).API
	var diags diag.Diagnostics

	// For logging purposes, only use the first four digits of the key,
	apiKeyTrunc := Truncate(d.Id(), 4)

	log.Printf("[DEBUG] Deleting API key [%s]", apiKeyTrunc)
	_, deleteErr := client.ApiKeys.Delete(d.Id())
	if deleteErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("unable to delete API key [%s]", apiKeyTrunc),
			Detail:   deleteErr.Error(),
		})
		return diags
	}

	log.Printf("[DEBUG] Deleted API key [%s]", apiKeyTrunc)

	d.SetId("")

	return diags
}
