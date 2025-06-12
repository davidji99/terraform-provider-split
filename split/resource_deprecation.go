package split

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// isHarnessTokenSet checks if the harness_token is set in the provider configuration
func isHarnessTokenSet(meta interface{}) bool {
	config := meta.(*Config)
	return config.harnessToken != ""
}

// checkResourceDeprecationWithHarnessToken returns diagnostics if the resource is deprecated when harness_token is set
func checkResourceDeprecationWithHarnessToken(resourceName string, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	if isHarnessTokenSet(meta) {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Resource %s cannot be used when harness_token is set", resourceName),
			Detail:   fmt.Sprintf("The resource %s is deprecated when using harness_token for authentication. Please use the Harness Terraform provider instead.", resourceName),
		})
	}

	return diags
}

// checkApiKeyTypeDeprecationWithHarnessToken returns diagnostics if the api_key type is "admin" and harness_token is set
func checkApiKeyTypeDeprecationWithHarnessToken(resourceName string, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	if isHarnessTokenSet(meta) {
		if keyType, ok := d.GetOk("type"); ok && keyType.(string) == "admin" {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Resource %s with type 'admin' cannot be used when harness_token is set", resourceName),
				Detail:   fmt.Sprintf("The resource %s with type 'admin' is deprecated when using harness_token for authentication. Please use the Harness Terraform provider instead.", resourceName),
			})
		}
	}

	return diags
}

// checkApiKeyTypeDeprecationWithHarnessTokenDiff performs the same validation as checkApiKeyTypeDeprecationWithHarnessToken but for ResourceDiff
func checkApiKeyTypeDeprecationWithHarnessTokenDiff(resourceName string, d *schema.ResourceDiff, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	if isHarnessTokenSet(meta) {
		if keyType, ok := d.GetOk("type"); ok && keyType.(string) == "admin" {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Resource %s with type 'admin' cannot be used when harness_token is set", resourceName),
				Detail:   fmt.Sprintf("The resource %s with type 'admin' is deprecated when using harness_token for authentication. Please use the Harness Terraform provider instead.", resourceName),
			})
		}
	}

	return diags
}
