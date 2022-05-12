package split

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"strings"
	"testing"
)

func TestAccSplitAttribute_Basic(t *testing.T) {
	workspaceID := testAccConfig.GetWorkspaceIDorSkip(t)
	name := fmt.Sprintf("tftest-%s", acctest.RandString(8))
	editedName := strings.ReplaceAll(name, "tftest", "edited")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSplitAttribute_basic(workspaceID, name, "true"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"split_environment.foobar", "workspace_id", workspaceID),
					resource.TestCheckResourceAttr(
						"split_environment.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"split_environment.foobar", "production", "true"),
				),
			},
			{
				Config: testAccCheckSplitAttribute_basic(workspaceID, editedName, "false"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"split_environment.foobar", "workspace_id", workspaceID),
					resource.TestCheckResourceAttr(
						"split_environment.foobar", "name", editedName),
					resource.TestCheckResourceAttr(
						"split_environment.foobar", "production", "false"),
				),
			},
		},
	})
}

func testAccCheckSplitAttribute_basic(workspaceID, name, production string) string {
	return fmt.Sprintf(`
provider "split" {
	remove_environment_from_state_only = true
}

resource "split_environment" "foobar" {
	workspace_id = "%s"
	name = "%s"
	production = %s
}
`, workspaceID, name, production)
}
