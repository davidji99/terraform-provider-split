package split

import (
	"os"
	"strings"
	"testing"
)

// testAccGetProviderConfig returns a provider configuration string with the appropriate
// authentication credentials based on what environment variables are set
func testAccGetProviderConfig(additionalConfig ...string) string {
	var lines []string

	// Start the provider block
	lines = append(lines, `provider "split" {`)

	// Add additional config if provided
	for _, config := range additionalConfig {
		if config != "" {
			lines = append(lines, "\t"+config)
		}
	}

	// Add authentication based on environment variables
	if harnessToken := os.Getenv("HARNESS_TOKEN"); harnessToken != "" {
		lines = append(lines, "\tharness_token = \""+harnessToken+"\"")
	} else if apiKey := os.Getenv("SPLIT_API_KEY"); apiKey != "" {
		lines = append(lines, "\tapi_key = \""+apiKey+"\"")
	}

	// Add common config that's needed for tests
	lines = append(lines, "\tremove_environment_from_state_only = true")

	// Close the provider block
	lines = append(lines, "}")

	return strings.Join(lines, "\n")
}

// isHarnessTokenEnvironmentSet checks if we're running tests with harness_token
func isHarnessTokenEnvironmentSet() bool {
	return os.Getenv("HARNESS_TOKEN") != ""
}

// skipIfUsingHarnessToken skips tests for resources that are deprecated when harness_token is used
func skipIfUsingHarnessToken(t *testing.T, resourceName string) {
	if isHarnessTokenEnvironmentSet() {
		t.Skipf("Skipping test for %s as it is deprecated when using harness_token", resourceName)
	}
}

// skipIfUsingHarnessTokenAndAdminType skips tests for api_key resources with type="admin" when harness_token is used
func skipIfUsingHarnessTokenAndAdminType(t *testing.T, keyType string) {
	if isHarnessTokenEnvironmentSet() && keyType == "admin" {
		t.Skipf("Skipping test for split_api_key with type 'admin' as it is deprecated when using harness_token")
	}
}
