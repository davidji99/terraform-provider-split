package split

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccSplitSegmentEnvironmentAssociation_Basic(t *testing.T) {
	workspaceID := testAccConfig.GetWorkspaceIDorSkip(t)
	envName := fmt.Sprintf("tftest-env-%s", acctest.RandString(3))
	trafficTypeName := fmt.Sprintf("tftest-tt-%s", acctest.RandString(8))
	segmentName := fmt.Sprintf("tftest-seg-%s", acctest.RandString(8))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSplitSegmentEnvironmentAssociation_basic(workspaceID,
					envName, trafficTypeName, segmentName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"split_segment_environment_association.foobar", "workspace_id", workspaceID),
					resource.TestCheckResourceAttr(
						"split_segment_environment_association.foobar", "segment_name", segmentName),
					resource.TestCheckResourceAttrSet(
						"split_segment_environment_association.foobar", "environment_id"),
				),
			},
		},
	})
}

func testAccCheckSplitSegmentEnvironmentAssociation_basic(workspaceID, envName, trafficTypeName, segmentName string) string {
	return fmt.Sprintf(`
provider "split" {
	remove_environment_from_state_only = true
}

resource "split_environment" "foobar" {
	workspace_id = "%[1]s"
	name = "%[2]s"
	production = false
}

resource "split_traffic_type" "foobar" {
	workspace_id = "%[1]s"
	name = "%[3]s"
}

resource "split_segment" "foobar" {
	workspace_id = "%[1]s"
	traffic_type_id = split_traffic_type.foobar.id
	name = "%[4]s"
	description = "made by TF tester"
}

resource "split_segment_environment_association" "foobar" {
	workspace_id = "%[1]s"
	environment_id = split_environment.foobar.id
	segment_name = split_segment.foobar.name
}
`, workspaceID, envName, trafficTypeName, segmentName)
}
