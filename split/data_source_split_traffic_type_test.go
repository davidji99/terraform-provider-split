package split

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccDatasourceSplitTrafficType_Basic(t *testing.T) {
	workspaceID := testAccConfig.GetWorkspaceIDorSkip(t)
	name := testAccConfig.GetTrafficTypeNameorSkip(t)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSplitTrafficTypeDataSource_Basic(workspaceID, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.split_traffic_type.foobar", "workspace_id", workspaceID),
					resource.TestCheckResourceAttr(
						"data.split_traffic_type.foobar", "name", name),
					resource.TestCheckResourceAttrSet(
						"data.split_traffic_type.foobar", "type"),
					resource.TestCheckResourceAttrSet(
						"data.split_traffic_type.foobar", "display_attribute_id"),
					resource.TestCheckResourceAttrSet(
						"data.split_traffic_type.foobar", "workspace_name"),
				),
			},
		},
	})
}

func testAccCheckSplitTrafficTypeDataSource_Basic(workspaceID, name string) string {
	return fmt.Sprintf(`
data "split_traffic_type" "foobar" {
  workspace_id = "%s"
  name = "%s"
}
`, workspaceID, name)
}
