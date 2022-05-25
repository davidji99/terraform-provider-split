package split

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccSplitTrafficType_Basic(t *testing.T) {
	workspaceID := testAccConfig.GetWorkspaceIDorSkip(t)
	name := fmt.Sprintf("tt-tftest-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSplitTrafficType_basic(workspaceID, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"split_traffic_type.foobar", "workspace_id", workspaceID),
					resource.TestCheckResourceAttr(
						"split_traffic_type.foobar", "name", name),
					resource.TestCheckResourceAttrSet(
						"split_traffic_type.foobar", "type"),
					resource.TestCheckResourceAttrSet(
						"split_traffic_type.foobar", "display_attribute_id"),
				),
			},
		},
	})
}

func testAccCheckSplitTrafficType_basic(workspaceID, name string) string {
	return fmt.Sprintf(`
provider "split" {
	remove_environment_from_state_only = true
}

resource "split_traffic_type" "foobar" {
	workspace_id = "%s"
	name = "%s"
}
`, workspaceID, name)
}
