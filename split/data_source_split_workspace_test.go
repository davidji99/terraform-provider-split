package split

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccDatasourceSplitWorkspace_Basic(t *testing.T) {
	name := testAccConfig.GetWorkspaceNameorSkip(t)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSplitWorkspaceDataSource_Basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.split_workspace.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"data.split_workspace.foobar", "type", "workspace"),
					resource.TestCheckResourceAttr(
						"data.split_workspace.foobar", "requires_title_and_comments", "false"),
				),
			},
		},
	})
}

func testAccCheckSplitWorkspaceDataSource_Basic(name string) string {
	return fmt.Sprintf(`
data "split_workspace" "foobar" {
  name = "%s"
}
`, name)
}
