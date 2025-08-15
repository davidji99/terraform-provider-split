package split

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccSplitFlagSetImport_Basic(t *testing.T) {
	workspaceID := testAccConfig.GetWorkspaceIDorSkip(t)
	name := fmt.Sprintf("tftest-import-%s", acctest.RandString(8))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSplitFlagSetImport_basic(workspaceID, name),
			},
			{
				ResourceName:      "split_flag_set.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckSplitFlagSetImport_basic(workspaceID, name string) string {
	return fmt.Sprintf(`
resource "split_flag_set" "foobar" {
	workspace_id = "%s"
	name = "%s"
}
`, workspaceID, name)
}
