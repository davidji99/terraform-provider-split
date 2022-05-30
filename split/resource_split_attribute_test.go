package split

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccSplitAttribute_Basic(t *testing.T) {
	workspaceID := testAccConfig.GetWorkspaceIDorSkip(t)
	ttName := fmt.Sprintf("tt-tftest-%s", acctest.RandString(8))
	attrIdentifier := acctest.RandString(8)
	attrName := fmt.Sprintf("attr-tftest-%s", acctest.RandString(8))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSplitAttribute_basic(workspaceID, attrIdentifier, ttName, attrName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"split_attribute.foobar", "workspace_id", workspaceID),
					resource.TestCheckResourceAttrSet(
						"split_attribute.foobar", "traffic_type_id"),
					resource.TestCheckResourceAttr(
						"split_attribute.foobar", "display_name", attrName),
					resource.TestCheckResourceAttr(
						"split_attribute.foobar", "description", "this is my attribute description"),
					resource.TestCheckResourceAttr(
						"split_attribute.foobar", "data_type", "string"),
					resource.TestCheckResourceAttr(
						"split_attribute.foobar", "is_searchable", "true"),
				),
			},
		},
	})
}

func testAccCheckSplitAttribute_basic(workspaceID, attrIdentifier, ttName, attrName string) string {
	return fmt.Sprintf(`
provider "split" {
	remove_environment_from_state_only = true
}

resource "split_traffic_type" "foobar" {
	workspace_id = "%[1]s"
	name = "%[3]s"
}

resource "split_attribute" "foobar" {
	workspace_id = "%[1]s"
	traffic_type_id = split_traffic_type.foobar.id
	identifier = "%[2]s"
	display_name = "%[4]s"
	description = "this is my attribute description"
	data_type = "STRING"
	is_searchable = true
}
`, workspaceID, attrIdentifier, ttName, attrName)
}
