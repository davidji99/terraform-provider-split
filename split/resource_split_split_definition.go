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

func resourceSplitSplitDefinition() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSplitSplitDefinitionCreate,
		ReadContext:   resourceSplitSplitDefinitionRead,
		UpdateContext: resourceSplitSplitDefinitionUpdate,
		DeleteContext: resourceSplitSplitDefinitionDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceSplitSplitDefinitionImport,
		},

		Schema: map[string]*schema.Schema{
			"workspace_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"split_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"environment_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"default_treatment": {
				Type:     schema.TypeString,
				Required: true,
			},

			"treatment": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"configurations": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsJSON,
						},

						"description": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},

			"default_rule": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"treatment": {
							Type:     schema.TypeString,
							Required: true,
						},

						"size": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},

			"rule": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bucket": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"treatment": {
										Type:     schema.TypeString,
										Required: true,
									},

									"size": {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
						},

						"condition": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"matcher": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:     schema.TypeString,
													Required: true,
												},

												"attribute": {
													Type:     schema.TypeString,
													Required: true,
												},

												"strings": {
													Type: schema.TypeList,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Required: true,
												},
											},
										},
									},

									"combiner": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceSplitSplitDefinitionImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*Config).API

	importID, parseErr := parseCompositeID(d.Id(), 3)
	if parseErr != nil {
		return nil, parseErr
	}

	workspaceID := importID[0]
	splitName := importID[1]
	environmentID := importID[2]

	sd, _, getErr := client.Splits.GetDefinition(workspaceID, splitName, environmentID)
	if getErr != nil {
		return nil, getErr
	}

	d.SetId(sd.GetID())
	d.Set("workspace_id", workspaceID)
	d.Set("split_name", sd.GetName())
	d.Set("environment_id", sd.GetEnvironment().GetID())
	d.Set("default_treatment", sd.GetDefaultTreatment())

	// Set Treatment in state
	setTreatmentInState(d, sd)

	// Set default rule
	setDefaultRuleInState(d, sd)

	// Set rule in state
	setRuleInState(d, sd)

	return []*schema.ResourceData{d}, nil
}

func resourceSplitSplitDefinitionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*Config).API

	opts := constructSplitDefinitionRequestOpts(d)
	workspaceID := getWorkspaceID(d)
	splitName := getSplitName(d)
	environmentID := getEnvironmentID(d)

	log.Printf("[DEBUG] Creating definition on split [%v]", splitName)

	sd, _, createErr := client.Splits.CreateDefinition(workspaceID, splitName, environmentID, opts)
	if createErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Unable to create definition for split %v", splitName),
			Detail:   createErr.Error(),
		})
		return diags
	}

	log.Printf("[DEBUG] Created definition on split [%v]", splitName)

	d.SetId(sd.GetID())

	return resourceSplitSplitDefinitionRead(ctx, d, meta)
}

func resourceSplitSplitDefinitionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*Config).API

	opts := constructSplitDefinitionRequestOpts(d)
	workspaceID := getWorkspaceID(d)
	environmentID := getEnvironmentID(d)
	splitName := getSplitName(d)

	log.Printf("[DEBUG] Updating split definition %v", d.Id())

	_, _, updateErr := client.Splits.UpdateDefinitionFull(workspaceID, splitName, environmentID, opts)
	if updateErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Unable to update split definition %v", d.Id()),
			Detail:   updateErr.Error(),
		})
		return diags
	}

	log.Printf("[DEBUG] Updating split definition %v", d.Id())

	return resourceSplitSplitDefinitionRead(ctx, d, meta)
}

func resourceSplitSplitDefinitionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*Config).API

	workspaceID := getWorkspaceID(d)
	environmentID := getEnvironmentID(d)
	splitName := getSplitName(d)

	sd, _, getErr := client.Splits.GetDefinition(workspaceID, splitName, environmentID)
	if getErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("unable to fetch split %s", d.Id()),
			Detail:   getErr.Error(),
		})
		return diags
	}

	d.Set("workspace_id", workspaceID)
	d.Set("split_name", sd.GetName())
	d.Set("environment_id", sd.GetEnvironment().GetID())
	d.Set("default_treatment", sd.GetDefaultTreatment())

	// Set Treatment in state
	setTreatmentInState(d, sd)

	// Set default rule
	setDefaultRuleInState(d, sd)

	// Set rule in state
	setRuleInState(d, sd)

	return diags
}

func resourceSplitSplitDefinitionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Config).API
	var diags diag.Diagnostics

	workspaceID := getWorkspaceID(d)
	envID := getEnvironmentID(d)
	splitName := getSplitName(d)

	log.Printf("[DEBUG] Deleting split definition %s", d.Id())

	_, deleteErr := client.Splits.RemoveDefinition(workspaceID, splitName, envID)
	if deleteErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("unable to delete split definition %s", d.Id()),
			Detail:   deleteErr.Error(),
		})
		return diags
	}

	log.Printf("[DEBUG] Deleted split defintion %s", d.Id())

	d.SetId("")

	return diags
}

func constructSplitDefinitionRequestOpts(d *schema.ResourceData) *api.SplitDefinitionRequest {
	opts := &api.SplitDefinitionRequest{}

	if v, ok := d.GetOk("default_treatment"); ok {
		opts.DefaultTreatment = v.(string)
		log.Printf("[DEBUG] new split definition default_treatment is : %v", opts.DefaultTreatment)
	}

	if v, ok := d.GetOk("treatment"); ok {
		treatments := make([]api.Treatment, 0)
		vL := v.([]interface{})

		for _, v := range vL {
			vt := v.(map[string]interface{})
			t := api.Treatment{}

			if v, ok := vt["name"].(string); ok {
				t.Name = &v
			}

			if v, ok := vt["configurations"].(string); ok {
				t.Configurations = &v
			}

			if v, ok := vt["description"].(string); ok {
				t.Description = &v
			}

			treatments = append(treatments, t)
		}

		opts.Treatments = treatments
		log.Printf("[DEBUG] new split definition treatments are : %v", opts.Treatments)
	}

	if v, ok := d.GetOk("default_rule"); ok {
		defaultRules := make([]api.Bucket, 0)
		vL := v.([]interface{})

		for _, v := range vL {
			var treatment string
			var size int
			vt := v.(map[string]interface{})

			if v, ok := vt["treatment"].(string); ok {
				treatment = v
			}

			if v, ok := vt["size"].(int); ok {
				size = v
			}
			defaultRules = append(defaultRules, api.Bucket{Treatment: &treatment, Size: &size})
		}

		opts.DefaultRule = defaultRules
		log.Printf("[DEBUG] new split definition default_rule is : %v", opts.DefaultRule)
	}

	if ruleRaw, ok := d.GetOk("rule"); ok {
		newRules := make([]api.Rule, 0)

		ruleListRaw := ruleRaw.([]interface{})

		for _, ruleRaw := range ruleListRaw {
			newRule := api.Rule{}
			rule := ruleRaw.(map[string]interface{})

			if bucketListRaw, ok := rule["bucket"].([]interface{}); ok {
				newRuleBuckets := make([]*api.Bucket, 0)
				for _, bucketRaw := range bucketListRaw {
					bucket := bucketRaw.(map[string]interface{})
					newRuleBucket := api.Bucket{}

					if v, ok := bucket["treatment"].(string); ok {
						newRuleBucket.Treatment = &v
					}
					if v, ok := bucket["size"].(int); ok {
						newRuleBucket.Size = &v
					}
					newRuleBuckets = append(newRuleBuckets, &newRuleBucket)
				}
				newRule.Buckets = newRuleBuckets
			}

			if conditionRaw, ok := rule["condition"]; ok {
				condition := conditionRaw.([]interface{})
				newRuleCondition := api.Condition{}

				for _, c := range condition {
					if combiner, ok := c.(map[string]interface{})["combiner"].(string); ok {
						newRuleCondition.Combiner = &combiner
					}

					if matchersListRaw, ok := c.(map[string]interface{})["matcher"]; ok {
						matchersList := matchersListRaw.([]interface{})
						newRuleConditionMatchers := make([]*api.Matcher, 0)
						for _, matcherRaw := range matchersList {
							matcher := matcherRaw.(map[string]interface{})
							newRuleConditionMatcher := api.Matcher{}

							if v, ok := matcher["type"].(string); ok {
								newRuleConditionMatcher.Type = &v
							}
							if v, ok := matcher["attribute"].(string); ok {
								newRuleConditionMatcher.Attribute = &v
							}
							if v, ok := matcher["strings"]; ok {
								stringsRaw := v.([]interface{})
								sList := make([]string, 0)
								for _, s := range stringsRaw {
									sList = append(sList, s.(string))
								}
								newRuleConditionMatcher.Strings = sList
							}
							newRuleConditionMatchers = append(newRuleConditionMatchers, &newRuleConditionMatcher)
						}
						newRuleCondition.Matchers = newRuleConditionMatchers
					}
				}
				newRule.Condition = &newRuleCondition
			}
			newRules = append(newRules, newRule)
		}
		opts.Rules = newRules
	}

	return opts
}

func setTreatmentInState(d *schema.ResourceData, sd *api.SplitDefinition) {
	treatments := make([]map[string]string, 0)
	for _, t := range sd.Treatments {
		treatments = append(treatments, map[string]string{
			"name":           t.GetName(),
			"configurations": t.GetConfigurations(),
			"description":    t.GetDescription(),
		})
	}
	d.Set("treatment", treatments)
}

func setDefaultRuleInState(d *schema.ResourceData, sd *api.SplitDefinition) {
	defaultRule := make([]map[string]interface{}, 0)
	for _, dr := range sd.DefaultRule {
		defaultRule = append(defaultRule, map[string]interface{}{
			"treatment": dr.GetTreatment(),
			"size":      dr.GetSize(),
		})
	}
	d.Set("default_rule", defaultRule)
}

func setRuleInState(d *schema.ResourceData, sd *api.SplitDefinition) {
	rules := make([]map[string]interface{}, 0)
	for _, r := range sd.Rules {
		buckets := make([]map[string]interface{}, 0)
		for _, b := range r.Buckets {
			buckets = append(buckets, map[string]interface{}{
				"treatment": b.GetTreatment(),
				"size":      b.GetSize(),
			})
		}

		conditions := make([]map[string]interface{}, 0)
		ruleConditionMatchers := make([]map[string]interface{}, 0)
		for _, rcm := range r.Condition.Matchers {
			ruleConditionMatchers = append(ruleConditionMatchers, map[string]interface{}{
				"type":      rcm.GetType(),
				"attribute": rcm.GetAttribute(),
				"strings":   rcm.Strings,
			})
		}
		conditions = append(conditions, map[string]interface{}{
			"combiner": *r.Condition.Combiner,
			"matcher":  ruleConditionMatchers,
		})

		rules = append(rules, map[string]interface{}{
			"bucket":    buckets,
			"condition": conditions,
		})
	}
	d.Set("rule", rules)
}
