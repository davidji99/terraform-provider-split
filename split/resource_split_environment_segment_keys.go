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

func resourceSplitEnvironmentSegmentKeys() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSplitEnvironmentSegmentKeysCreate,
		ReadContext:   resourceSplitEnvironmentSegmentKeysRead,
		UpdateContext: resourceSplitEnvironmentSegmentKeysUpdate,
		DeleteContext: resourceSplitEnvironmentSegmentKeysDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceSplitEnvironmentSegmentKeysImport,
		},

		Schema: map[string]*schema.Schema{
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

			"keys": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
				MinItems: 1,
				MaxItems: 10000,
			},
		},
	}
}

func resourceSplitEnvironmentSegmentKeysImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*Config).API

	importID, parseErr := parseCompositeID(d.Id(), 2)
	if parseErr != nil {
		return nil, parseErr
	}

	environmentID := importID[0]
	segmentName := importID[1]

	result, _, getErr := client.Environments.GetSegmentKeys(environmentID, segmentName)
	if getErr != nil {
		return nil, fmt.Errorf(fmt.Sprintf("unable to fetch environment %s segment %s's keys", environmentID, segmentName))
	}

	d.SetId(fmt.Sprintf("%s:%s", environmentID, segmentName))

	d.Set("environment_id", environmentID)
	d.Set("segment_name", segmentName)

	keys := make([]string, 0)
	for _, key := range result.Keys {
		keys = append(keys, key.GetKey())
	}
	d.Set("keys", keys)

	return []*schema.ResourceData{d}, nil
}

func resourceSplitEnvironmentSegmentKeysCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*Config).API
	opts := &api.EnvironmentSegmentKeysRequest{}
	environmentID := getEnvironmentID(d)
	segmentName := d.Get("segment_name").(string)

	if v, ok := d.GetOk("keys"); ok {
		vL := v.(*schema.Set).List()
		keys := make([]string, 0)
		for _, e := range vL {
			keys = append(keys, e.(string))
		}
		opts.Keys = keys
		log.Printf("[DEBUG] adding keys : %v", opts.Keys)
	}

	log.Printf("[DEBUG] Modifying segment keys to environment %s & segment %s", environmentID, segmentName)

	_, _, addErr := client.Environments.AddSegmentKeys(environmentID, segmentName, true, opts)
	if addErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("[DEBUG] Unable to add segment keys to environment %s & segment %s", environmentID, segmentName),
			Detail:   addErr.Error(),
		})
		return diags
	}

	log.Printf("[DEBUG] Added segment keys to environment %s & segment %s", environmentID, segmentName)

	d.SetId(fmt.Sprintf("%s:%s", environmentID, segmentName))

	return resourceSplitEnvironmentSegmentKeysRead(ctx, d, meta)
}

func resourceSplitEnvironmentSegmentKeysUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*Config).API
	environmentID := getEnvironmentID(d)
	segmentName := d.Get("segment_name").(string)

	hasChange := d.HasChange("keys")
	log.Printf("[INFO] Does segment environment association have changes: *%#v", hasChange)

	var oldKeysRaw, newKeysRaw []interface{}
	if hasChange {
		o, n := d.GetChange("keys")
		if o == nil {
			o = []interface{}{}
		}
		if n == nil {
			n = []interface{}{}
		}

		oldKeysRaw = o.(*schema.Set).List()
		newKeysRaw = n.(*schema.Set).List()
	}

	// First, remove the old keys
	if len(oldKeysRaw) > 0 {
		oldKeys := make([]string, 0, len(oldKeysRaw))
		for _, k := range oldKeysRaw {
			oldKeys = append(oldKeys, k.(string))
		}

		log.Printf("[INFO] Removing the following keys from environment [%s] and segment [%s]: %v", environmentID, segmentName, oldKeys)

		opts := &api.EnvironmentSegmentKeysRequest{
			Keys:    oldKeys,
			Comment: "modified by Terraform",
		}
		_, err := client.Environments.RemoveSegmentKeys(environmentID, segmentName, opts)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Unable to remove keys when updating environment [%s] and segment [%s]", environmentID, segmentName),
				Detail:   err.Error(),
			})
			return diags
		}

		log.Printf("[INFO] Removed the following keys from environment [%s] and segment [%s]: %v", environmentID, segmentName, oldKeys)
	}

	// Then, add the new keys
	if len(newKeysRaw) > 0 {
		newKeys := make([]string, 0, len(newKeysRaw))
		for _, k := range newKeysRaw {
			newKeys = append(newKeys, k.(string))
		}

		log.Printf("[INFO] Adding the following keys from environment [%s] and segment [%s]: %v", environmentID, segmentName, newKeys)

		opts := &api.EnvironmentSegmentKeysRequest{
			Keys:    newKeys,
			Comment: "modified by Terraform",
		}
		_, _, err := client.Environments.AddSegmentKeys(environmentID, segmentName, true, opts)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Unable to add keys when updating environment [%s] and segment [%s]", environmentID, segmentName),
				Detail:   err.Error(),
			})
			return diags
		}

		log.Printf("[INFO] Added the following keys from environment [%s] and segment [%s]: %v", environmentID, segmentName, newKeys)
	}

	return resourceSplitEnvironmentSegmentKeysRead(ctx, d, meta)
}

func resourceSplitEnvironmentSegmentKeysRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*Config).API
	environmentID := getEnvironmentID(d)
	segmentName := d.Get("segment_name").(string)

	result, _, getErr := client.Environments.GetSegmentKeys(environmentID, segmentName)
	if getErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("unable to fetch environment %s segment %s's keys", environmentID, segmentName),
			Detail:   getErr.Error(),
		})
		return diags
	}

	d.Set("environment_id", environmentID)
	d.Set("segment_name", segmentName)

	keys := make([]string, 0)
	for _, key := range result.Keys {
		keys = append(keys, key.GetKey())
	}
	d.Set("keys", keys)

	return diags
}

func resourceSplitEnvironmentSegmentKeysDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*Config).API

	result, parseErr := parseCompositeID(d.Id(), 2)
	if parseErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "unable to parse resource ID during deletion",
			Detail:   parseErr.Error(),
		})
		return diags
	}
	environmentId := result[0]
	segmentName := result[1]

	log.Printf("[DEBUG] Removing all segment keys from environment %s & segment %s due to resource deletion", environmentId, segmentName)

	keysRaw := d.Get("keys").(*schema.Set).List()
	keys := make([]string, len(keysRaw))
	for _, k := range keysRaw {
		keys = append(keys, k.(string))
	}

	opts := &api.EnvironmentSegmentKeysRequest{
		Keys:    keys,
		Comment: "modified by Terraform",
	}

	_, deleteErr := client.Environments.RemoveSegmentKeys(environmentId, segmentName, opts)
	if deleteErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("unable to delete all keys from environment %s segment %s's keys", environmentId, segmentName),
			Detail:   deleteErr.Error(),
		})
		return diags
	}

	log.Printf("[DEBUG] Removed all segment keys from environment %s & segment %s due to resource deletion", environmentId, segmentName)

	d.SetId("")

	return diags
}
