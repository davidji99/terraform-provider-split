package split

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccSplitSegment_Basic(t *testing.T) {
	workspaceID := testAccConfig.GetWorkspaceIDorSkip(t)
	trafficTypeID := testAccConfig.GetTrafficTypeNameorSkip(t)
	name := fmt.Sprintf("tftest-%s", acctest.RandString(8))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSplitSegment_basic(workspaceID, trafficTypeID, name, "created from Terraform"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"split_segment.foobar", "workspace_id", workspaceID),
					resource.TestCheckResourceAttr(
						"split_segment.foobar", "traffic_type_id", trafficTypeID),
					resource.TestCheckResourceAttr(
						"split_segment.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"split_segment.foobar", "description", "created from Terraform"),
				),
			},
		},
	})
}

func testAccCheckSplitSegment_basic(workspaceID, trafficTypeID, name, description string) string {
	return fmt.Sprintf(`
resource "split_segment" "foobar" {
	workspace_id = "%s"
	traffic_type_id = "%s"
	name = "%s"
	description = "%s"
}
`, workspaceID, trafficTypeID, name, description)
}
