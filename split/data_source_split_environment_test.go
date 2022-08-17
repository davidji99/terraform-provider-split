package split

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccDatasourceSplitEnvironment_Basic(t *testing.T) {
	workspaceID := testAccConfig.GetWorkspaceIDorSkip(t)
	name := fmt.Sprintf("tftest-%s", acctest.RandString(8))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSplitEnvironmentDataSource_Basic(workspaceID, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.split_environment.foobar", "workspace_id", workspaceID),
					resource.TestCheckResourceAttr(
						"data.split_environment.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"data.split_environment.foobar", "production", "false"),
				),
			},
		},
	})
}

func testAccCheckSplitEnvironmentDataSource_Basic(workspaceID, envName string) string {
	return fmt.Sprintf(`
provider "split" {
	remove_environment_from_state_only = true
}

resource "split_environment" "foobar" {
	workspace_id = "%[1]s"
	name = "%[2]s"
}

data "split_environment" "foobar" {
  workspace_id = "%[1]s"
  name = split_environment.foobar.name
}
`, workspaceID, envName)
}
