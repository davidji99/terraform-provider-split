package split

import (
	"testing"

	helper "github.com/davidji99/terraform-provider-split/helper/test"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var providers []*schema.Provider
var testAccProviderFactories map[string]func() (*schema.Provider, error)
var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider
var testAccConfig *helper.TestConfig

func init() {
	testAccProvider = New()
	testAccProviders = map[string]*schema.Provider{
		"split": testAccProvider,
	}
	testAccProviderFactories = testAccProviderFactoriesInit(providers, []string{"split"})
	testAccConfig = helper.NewTestConfig()
}

func TestProvider(t *testing.T) {
	if err := New().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = New()
}

func testAccPreCheck(t *testing.T) {
	testAccConfig.GetOrAbort(t, helper.TestConfigSplitAPIKey)
}

func testAccProviderFactoriesInit(providers []*schema.Provider, providerNames []string) map[string]func() (*schema.Provider, error) {
	var factories = make(map[string]func() (*schema.Provider, error), len(providerNames))

	for _, name := range providerNames {
		p := New()

		factories[name] = func() (*schema.Provider, error) {
			return p, nil
		}

		if providers != nil {
			providers = append(providers, p)
		}
	}

	return factories
}
