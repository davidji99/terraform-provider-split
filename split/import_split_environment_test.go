package split

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
)

func TestAccSplitEnvironment_importBasic(t *testing.T) {
	workspaceID := testAccConfig.GetWorkspaceIDorSkip(t)
	name := fmt.Sprintf("tftest-%s", acctest.RandString(8))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSplitEnvironment_basic(workspaceID, name, "true"),
			},
			{
				ResourceName:      "split_environment.foobar",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccSplitEnvironmentImportStateIDFunc("split_environment.foobar"),
			},
		},
	})
}

func testAccSplitEnvironmentImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("[ERROR] Not found: %s", resourceName)
		}

		return fmt.Sprintf("%s:%s", rs.Primary.Attributes["workspace_id"], rs.Primary.Attributes["name"]), nil
	}
}
