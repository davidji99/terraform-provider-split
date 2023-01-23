package split

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
)

func TestAccSplitEnvironmentSegmentKeys_importBasic(t *testing.T) {
	workspaceID := testAccConfig.GetWorkspaceIDorSkip(t)
	envName := fmt.Sprintf("tftest-env-%s", acctest.RandString(3))
	trafficTypeName := fmt.Sprintf("tftest-tt-%s", acctest.RandString(8))
	segmentName := fmt.Sprintf("tftest-seg-%s", acctest.RandString(8))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSplitEnvironmentSegmentKeys_basic(workspaceID, envName, trafficTypeName,
					segmentName, true),
			},
			{
				ResourceName:      "split_environment_segment_keys.foobar",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccSplitEnvironmentSegmentKeysImportStateIDFunc("split_environment_segment_keys.foobar"),
			},
		},
	})
}

func testAccSplitEnvironmentSegmentKeysImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("[ERROR] Not found: %s", resourceName)
		}

		return fmt.Sprintf("%s:%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["segment_name"]), nil
	}
}
