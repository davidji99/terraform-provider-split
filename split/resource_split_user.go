package split

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pmcjury/terraform-provider-split/api"
)

func resourceSplitUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSplitUserCreate,
		ReadContext:   resourceSplitUserRead,
		UpdateContext: resourceSplitUserUpdate,
		DeleteContext: resourceSplitUserDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},

			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"2fa": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSplitUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*Config).API
	opts := &api.UserCreateRequest{}

	if v, ok := d.GetOk("email"); ok {
		opts.Email = v.(string)
		log.Printf("[DEBUG] new user email is : %v", opts.Email)
	}

	log.Printf("[DEBUG] Inviting user %s", opts.Email)

	u, _, inviteErr := client.Users.Invite(opts)
	if inviteErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Unable to create/invite user %v", opts.Email),
			Detail:   inviteErr.Error(),
		})
		return diags
	}

	log.Printf("[DEBUG] Invited user %s", opts.Email)

	d.SetId(u.GetID())

	return resourceSplitUserRead(ctx, d, meta)
}

func resourceSplitUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*Config).API

	u, _, getErr := client.Users.Get(d.Id())
	if getErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Unable to fetch user %v", d.Id()),
			Detail:   getErr.Error(),
		})
		return diags
	}

	d.Set("email", u.GetEmail())
	d.Set("name", u.GetName())
	d.Set("2fa", u.GetTFA())
	d.Set("status", u.GetStatus())

	return diags
}

func resourceSplitUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*Config).API
	opts := &api.UserUpdateRequest{}

	if v, ok := d.GetOk("name"); ok {
		opts.Name = v.(string)
		log.Printf("[DEBUG] updated user name is : %v", opts.Name)
	}

	_, _, updateErr := client.Users.Update(d.Id(), opts)
	if updateErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Unable to update user %v", opts.Email),
			Detail:   updateErr.Error(),
		})
		return diags
	}

	return resourceSplitUserRead(ctx, d, meta)
}

func resourceSplitUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*Config).API

	// Check the status of the user prior to deletion.
	u, _, getErr := client.Users.Get(d.Id())
	if getErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Unable to fetch user %v for deletion", d.Id()),
			Detail:   getErr.Error(),
		})
		return diags
	}

	// If the user's status is 'PENDING', delete the invitation.
	if u.GetStatus() == api.UserStatusPending {
		log.Printf("[DEBUG] Deleting invitation for user %s", d.Id())

		_, deleteErr := client.Users.DeletePendingUser(d.Id())
		if deleteErr != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Unable to delete invitation for user %v", d.Id()),
				Detail:   deleteErr.Error(),
			})
			return diags
		}

		log.Printf("[DEBUG] Deleted invitation for user %s", d.Id())
	}

	// If the user's status is 'ACTIVE', deactivate the user.
	if u.GetStatus() == api.UserStatusActive {
		log.Printf("[DEBUG] Disabling user %s", d.Id())

		_, _, deleteErr := client.Users.Update(d.Id(), &api.UserUpdateRequest{Status: api.UserStatusDeactivated})
		if deleteErr != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Unable to disable user %v", d.Id()),
				Detail:   deleteErr.Error(),
			})
			return diags
		}

		log.Printf("[DEBUG] Disabled user %s", d.Id())
	}

	if u.GetStatus() != api.UserStatusActive && u.GetStatus() != api.UserStatusPending {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Unable to disable/delete user %v", d.Id()),
			Detail:   fmt.Sprintf("unsupported user status: %s. Expected 'PENDING' or 'ACTIVE'", u.GetStatus()),
		})
		return diags
	}

	d.SetId("")

	return diags
}
