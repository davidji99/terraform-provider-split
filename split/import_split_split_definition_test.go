package split

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
)

func TestAccSplitSplitDefinition_importBasic(t *testing.T) {
	workspaceID := testAccConfig.GetWorkspaceIDorSkip(t)
	envID := testAccConfig.GetEnvironmentIDorSkip(t)
	trafficTypeID := testAccConfig.GetTrafficTypeIDorSkip(t)
	trafficTypeName := fmt.Sprintf("tt-tftest-%s", acctest.RandString(10))
	splitName := fmt.Sprintf("s-tftest-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSplitSplitDefinition_basicUpdated(workspaceID, trafficTypeName, splitName,
					"my split description", envID, trafficTypeID),
			},
			{
				ResourceName:      "split_split_definition.foobar",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccSplitSplitDefinitionImportStateIDFunc("split_split_definition.foobar"),
			},
		},
	})
}

func testAccSplitSplitDefinitionImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("[ERROR] Not found: %s", resourceName)
		}

		return fmt.Sprintf("%s:%s:%s", rs.Primary.Attributes["workspace_id"], rs.Primary.Attributes["id"],
			rs.Primary.Attributes["environment_id"]), nil
	}
}
