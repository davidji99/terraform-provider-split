package split

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
)

func TestAccSplitAttribute_importBasic(t *testing.T) {
	workspaceID := testAccConfig.GetWorkspaceIDorSkip(t)
	ttName := fmt.Sprintf("tt-tftest-%s", acctest.RandString(8))
	attrIdentifier := acctest.RandString(8)
	attrName := fmt.Sprintf("attr-tftest-%s", acctest.RandString(8))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSplitAttribute_basic(workspaceID, attrIdentifier, ttName, attrName),
			},
			{
				ResourceName:      "split_attribute.foobar",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccSplitAttributeImportStateIDFunc("split_attribute.foobar"),
			},
		},
	})
}

func testAccSplitAttributeImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("[ERROR] Not found: %s", resourceName)
		}

		return fmt.Sprintf("%s:%s:%s", rs.Primary.Attributes["workspace_id"], rs.Primary.Attributes["traffic_type_id"],
			rs.Primary.Attributes["id"]), nil
	}
}
