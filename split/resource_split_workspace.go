package split

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pmcjury/terraform-provider-split/api"
)

func resourceSplitWorkspace() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSplitWorkspaceCreate,
		ReadContext:   resourceSplitWorkspaceRead,
		UpdateContext: resourceSplitWorkspaceUpdate,
		DeleteContext: resourceSplitWorkspaceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceSplitWorkspaceImport,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"require_title_comments": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceSplitWorkspaceImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*Config).API

	w, _, getErr := client.Workspaces.FindByName(d.Id())
	if getErr != nil {
		return nil, getErr
	}

	d.SetId(w.GetID())
	d.Set("name", w.GetName())
	d.Set("require_title_comments", w.GetRequiresTitleAndComments())

	return []*schema.ResourceData{d}, nil
}

func resourceSplitWorkspaceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*Config).API
	opts := &api.WorkspaceRequest{}

	if v, ok := d.GetOk("name"); ok {
		vs := v.(string)
		opts.Name = &vs
		log.Printf("[DEBUG] new workspace name is : %s", *opts.Name)
	}

	rtc := d.Get("require_title_comments").(bool)
	opts.RequiresTitleAndComments = &rtc
	log.Printf("[DEBUG] new workspace require_title_comments is : %v", *opts.RequiresTitleAndComments)

	log.Printf("[DEBUG] Creating new workspace %v", opts.Name)

	w, _, createErr := client.Workspaces.Create(opts)
	if createErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Unable to create workspace %v", opts.Name),
			Detail:   createErr.Error(),
		})
		return diags
	}

	log.Printf("[DEBUG] Created new workspace %v", w.GetName())

	d.SetId(w.GetID())
	d.Set("name", w.GetName())

	return resourceSplitWorkspaceRead(ctx, d, meta)
}

func resourceSplitWorkspaceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*Config).API

	opts := &api.WorkspaceRequest{}

	if ok := d.HasChange("name"); ok {
		vs := d.Get("name").(string)
		opts.Name = &vs
		log.Printf("[DEBUG] updated workspace name is : %s", *opts.Name)
	}

	if ok := d.HasChange("require_title_comments"); ok {
		vs := d.Get("require_title_comments").(bool)
		opts.RequiresTitleAndComments = &vs
		log.Printf("[DEBUG] updated workspace require_title_comments is : %v", *opts.RequiresTitleAndComments)
	}

	log.Printf("[DEBUG] Updating workspace %v", d.Id())

	_, _, updateErr := client.Workspaces.Update(d.Id(), opts)
	if updateErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Unable to update workspace %v", d.Id()),
			Detail:   updateErr.Error(),
		})
		return diags
	}

	log.Printf("[DEBUG] Updated workspace %v", d.Id())

	return resourceSplitWorkspaceRead(ctx, d, meta)
}

func resourceSplitWorkspaceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*Config).API

	workspace, _, findErr := client.Workspaces.FindByName(d.Get("name").(string))
	if findErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("unable to fetch workspace %s", d.Id()),
			Detail:   findErr.Error(),
		})
		return diags
	}

	d.Set("name", workspace.GetName())
	d.Set("require_title_comments", workspace.GetRequiresTitleAndComments())

	return diags
}

func resourceSplitWorkspaceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Config).API
	var diags diag.Diagnostics

	// You can't delete a workspace if the workspace has child traffic types. By default, a `user` traffic type
	// is created with a new workspace. Therefore, we need to first delete the workspace traffic types before
	// deleting the workspace itself.
	log.Printf("[DEBUG] Finding traffic types prior to workspace %s deletion", d.Id())

	trafficTypes, _, listErr := client.TrafficTypes.List(d.Id())
	if listErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("unable to lookup traffic types for workspace %s", d.Id()),
			Detail:   listErr.Error(),
		})
		return diags
	}

	log.Printf("[DEBUG] Deleting all traffic types associated with workspace %s", d.Id())
	for _, tt := range trafficTypes {
		log.Printf("[DEBUG] Deleting traffic type %s associated with workspace %s", tt.GetID(), d.Id())
		_, deleteErr := client.TrafficTypes.Delete(tt.GetID())
		if deleteErr != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("unable to delete traffic type %s", d.Id()),
				Detail:   deleteErr.Error(),
			})
			return diags
		}
	}
	log.Printf("[DEBUG] Deleted all traffic types associated with workspace %s", d.Id())

	log.Printf("[DEBUG] Deleting workspace %s", d.Id())
	_, deleteErr := client.Workspaces.Delete(d.Id())
	if deleteErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("unable to delete workspace %s", d.Id()),
			Detail:   deleteErr.Error(),
		})
		return diags
	}

	log.Printf("[DEBUG] Deleted workspace %s", d.Id())

	d.SetId("")

	return diags
}
