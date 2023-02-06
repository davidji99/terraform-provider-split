package split

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccSplitEnvironmentSegmentKeys_Basic(t *testing.T) {
	workspaceID := testAccConfig.GetWorkspaceIDorSkip(t)
	envName := fmt.Sprintf("tftest-env-%s", acctest.RandString(3))
	trafficTypeName := fmt.Sprintf("tftest-tt-%s", acctest.RandString(8))
	segmentName := fmt.Sprintf("tftest-seg-%s", acctest.RandString(8))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSplitEnvironmentSegmentKeys_basic(workspaceID, envName, trafficTypeName,
					segmentName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"split_environment_segment_keys.foobar", "segment_name", segmentName),
					resource.TestCheckResourceAttr(
						"split_environment_segment_keys.foobar", "keys.#", "2"),
					//resource.TestCheckResourceAttr(
					//	"split_environment_segment_keys.foobar", "keys.#.0", "tester1"),
					//resource.TestCheckResourceAttr(
					//	"split_environment_segment_keys.foobar", "keys.#.1", "tester2"),
					resource.TestCheckResourceAttrSet(
						"split_environment_segment_keys.foobar", "environment_id"),
				),
			},
			{
				Config: testAccCheckSplitEnvironmentSegmentKeys_updatedKeys(workspaceID, envName, trafficTypeName,
					segmentName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"split_environment_segment_keys.foobar", "segment_name", segmentName),
					resource.TestCheckResourceAttr(
						"split_environment_segment_keys.foobar", "keys.#", "3"),
					//resource.TestCheckResourceAttr(
					//	"split_environment_segment_keys.foobar", "keys.#.0", "tester2"),
					//resource.TestCheckResourceAttr(
					//	"split_environment_segment_keys.foobar", "keys.#.1", "tester3"),
					//resource.TestCheckResourceAttr(
					//	"split_environment_segment_keys.foobar", "keys.#.2", "tester3"),
					resource.TestCheckResourceAttrSet(
						"split_environment_segment_keys.foobar", "environment_id"),
				),
			},
		},
	})
}

func testAccCheckSplitEnvironmentSegmentKeys_basic(
	workspaceID, environmentName, trafficTypeName, segmentName string, production bool) string {
	return fmt.Sprintf(`
provider "split" {
	remove_environment_from_state_only = true
}

resource "split_environment" "foobar" {
	workspace_id = "%[1]s"
	name = "%[2]s"
	production = %[5]v
}

resource "split_traffic_type" "foobar" {
	workspace_id = "%[1]s"
	name = "%[3]s"
}

resource "split_segment" "foobar" {
	workspace_id = "%[1]s"
	traffic_type_id = split_traffic_type.foobar.id
	name = "%[4]s"
	description = "description_of_my_segment"
}

resource "split_segment_environment_association" "foobar" {
	workspace_id = "%[1]s"
	environment_id = split_environment.foobar.id
	segment_name = split_segment.foobar.name
}

resource "split_environment_segment_keys" "foobar" {
	environment_id = split_environment.foobar.id
	segment_name = split_segment_environment_association.foobar.segment_name
	keys = ["tester1", "tester2"]
}

`, workspaceID, environmentName, trafficTypeName, segmentName, production)
}

func testAccCheckSplitEnvironmentSegmentKeys_updatedKeys(
	workspaceID, environmentName, trafficTypeName, segmentName string, production bool) string {
	return fmt.Sprintf(`
provider "split" {
	remove_environment_from_state_only = true
}

resource "split_environment" "foobar" {
	workspace_id = "%[1]s"
	name = "%[2]s"
	production = %[5]v
}

resource "split_traffic_type" "foobar" {
	workspace_id = "%[1]s"
	name = "%[3]s"
}

resource "split_segment" "foobar" {
	workspace_id = "%[1]s"
	traffic_type_id = split_traffic_type.foobar.id
	name = "%[4]s"
	description = "description_of_my_segment"
}

resource "split_segment_environment_association" "foobar" {
	workspace_id = "%[1]s"
	environment_id = split_environment.foobar.id
	segment_name = split_segment.foobar.name
}

resource "split_environment_segment_keys" "foobar" {
	environment_id = split_environment.foobar.id
	segment_name = split_segment_environment_association.foobar.segment_name
	keys = ["tester2", "tester3", "tester4"]
}

`, workspaceID, environmentName, trafficTypeName, segmentName, production)
}
