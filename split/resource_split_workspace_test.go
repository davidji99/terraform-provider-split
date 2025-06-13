package split

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccSplitWorkspace_Basic(t *testing.T) {
	// Skip test if using harness_token as this resource is deprecated with harness_token
	skipIfUsingHarnessToken(t, "split_workspace")

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
%s

resource "split_workspace" "foobar" {
	name = "%s"
	require_title_comments = %s
}
`, testAccGetProviderConfig(), name, requireTitleComments)
}
