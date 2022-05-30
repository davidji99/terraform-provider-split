package split

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccSplitWorkspace_Basic(t *testing.T) {
	name := fmt.Sprintf("w-tftest-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSplitWorkspace_basic(name, "true"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"split_workspace.foobar", "require_title_comments", "true"),
					resource.TestCheckResourceAttr(
						"split_workspace.foobar", "name", name),
				),
			},
			{
				Config: testAccCheckSplitWorkspace_basic(name+"edited", "false"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"split_workspace.foobar", "require_title_comments", "false"),
					resource.TestCheckResourceAttr(
						"split_workspace.foobar", "name", name+"edited"),
				),
			},
		},
	})
}

func testAccCheckSplitWorkspace_basic(name, requireTitleComments string) string {
	return fmt.Sprintf(`
provider "split" {
	remove_environment_from_state_only = true
}

resource "split_workspace" "foobar" {
	name = "%s"
	require_title_comments = %s
}
`, name, requireTitleComments)
}
