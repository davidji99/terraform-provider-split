package split

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccSplitApiKey_ClientSide_Basic(t *testing.T) {
	workspaceID := testAccConfig.GetWorkspaceIDorSkip(t)
	name := fmt.Sprintf("tftest-%s", acctest.RandString(8))
	envName := fmt.Sprintf("tftest-%s", acctest.RandString(8))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSplitApiKey_basic(workspaceID, name, "client_side", envName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"split_environment.foobar", "id"),
				),
			},
		},
	})
}

func TestAccSplitApiKey_ServerSide_Basic(t *testing.T) {
	workspaceID := testAccConfig.GetWorkspaceIDorSkip(t)
	name := fmt.Sprintf("tftest-%s", acctest.RandString(8))
	envName := fmt.Sprintf("tftest-%s", acctest.RandString(8))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSplitApiKey_basic(workspaceID, name, "server_side", envName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"split_environment.foobar", "id"),
				),
			},
		},
	})
}

// TestAccSplitApiKey_Admin_Basic may or may not work depending on your account's functionality. You may get the following
// error: `402 {"code":402,"message":"Forbidden by Paywalls:`
func TestAccSplitApiKey_Admin_Basic(t *testing.T) {
	// Skip test if using harness_token as admin API keys are deprecated with harness_token
	skipIfUsingHarnessTokenAndAdminType(t, "admin")

	workspaceID := testAccConfig.GetWorkspaceIDorSkip(t)
	name := fmt.Sprintf("tftest-%s", acctest.RandString(8))
	envName := fmt.Sprintf("tftest-%s", acctest.RandString(8))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSplitApiKey_basic_WithRoles(workspaceID, name, "admin", envName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"split_environment.foobar", "id"),
				),
			},
		},
	})
}

func testAccCheckSplitApiKey_basic(workspaceID, name, keyType, environmentName string) string {
	return fmt.Sprintf(`
provider "split" {
	remove_environment_from_state_only = true
}

resource "split_environment" "foobar" {
	workspace_id = "%[1]s"
	name = "%[4]s"
	production = true
}

resource "split_api_key" "foobar" {
	workspace_id = "%[1]s"
	name = "%[2]s"
	type = "%[3]s"
	environment_ids = [split_environment.foobar.id]
}
`, workspaceID, name, keyType, environmentName)
}

func testAccCheckSplitApiKey_basic_WithRoles(workspaceID, name, keyType, environmentName string) string {
	return fmt.Sprintf(`
provider "split" {
	remove_environment_from_state_only = true
}

resource "split_environment" "foobar" {
	workspace_id = "%[1]s"
	name = "%[4]s"
	production = true
}

resource "split_api_key" "foobar" {
	workspace_id = "%[1]s"
	name = "%[2]s"
	type = "%[3]s"
	environment_ids = [split_environment.foobar.id]
	roles = ["API_ALL_GRANTED"]
}
`, workspaceID, name, keyType, environmentName)
}
