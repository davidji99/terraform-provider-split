package split

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccSplitTrafficTypeAttribute_Basic(t *testing.T) {
	workspaceID := testAccConfig.GetWorkspaceIDorSkip(t)
	ttName := fmt.Sprintf("tt-tftest-%s", acctest.RandString(8))
	attrIdentifier := acctest.RandString(8)
	attrName := fmt.Sprintf("attr-tftest-%s", acctest.RandString(8))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSplitTrafficTypeAttribute_basic(workspaceID, attrIdentifier, ttName, attrName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"split_traffic_type_attribute.foobar", "workspace_id", workspaceID),
					resource.TestCheckResourceAttrSet(
						"split_traffic_type_attribute.foobar", "traffic_type_id"),
					resource.TestCheckResourceAttr(
						"split_traffic_type_attribute.foobar", "display_name", attrName),
					resource.TestCheckResourceAttr(
						"split_traffic_type_attribute.foobar", "description", "this is my attribute description"),
					resource.TestCheckResourceAttr(
						"split_traffic_type_attribute.foobar", "data_type", "STRING"),
					resource.TestCheckResourceAttr(
						"split_traffic_type_attribute.foobar", "suggested_values.#", "3"),
					resource.TestCheckResourceAttr(
						"split_traffic_type_attribute.foobar", "is_searchable", "true"),
					resource.TestCheckResourceAttrSet(
						"split_traffic_type_attribute.foobar", "organization_id"),
				),
			},
		},
	})
}

func TestAccSplitTrafficTypeAttribute_Updates(t *testing.T) {
	workspaceID := testAccConfig.GetWorkspaceIDorSkip(t)
	ttName := fmt.Sprintf("tt-tftest-%s", acctest.RandString(8))
	attrIdentifier := acctest.RandString(8)
	attrName := fmt.Sprintf("attr-tftest-%s", acctest.RandString(8))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSplitTrafficTypeAttribute_basic(workspaceID, attrIdentifier, ttName, attrName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"split_traffic_type_attribute.foobar", "workspace_id", workspaceID),
					resource.TestCheckResourceAttrSet(
						"split_traffic_type_attribute.foobar", "traffic_type_id"),
					resource.TestCheckResourceAttr(
						"split_traffic_type_attribute.foobar", "display_name", attrName),
					resource.TestCheckResourceAttr(
						"split_traffic_type_attribute.foobar", "description", "this is my attribute description"),
					resource.TestCheckResourceAttr(
						"split_traffic_type_attribute.foobar", "data_type", "STRING"),
					resource.TestCheckResourceAttr(
						"split_traffic_type_attribute.foobar", "suggested_values.#", "3"),
					resource.TestCheckResourceAttr(
						"split_traffic_type_attribute.foobar", "is_searchable", "true"),
					resource.TestCheckResourceAttrSet(
						"split_traffic_type_attribute.foobar", "organization_id"),
				),
			},
			{
				Config: testAccCheckSplitTrafficTypeAttribute_updates(workspaceID, attrIdentifier, ttName, attrName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"split_traffic_type_attribute.foobar", "workspace_id", workspaceID),
					resource.TestCheckResourceAttrSet(
						"split_traffic_type_attribute.foobar", "traffic_type_id"),
					resource.TestCheckResourceAttr(
						"split_traffic_type_attribute.foobar", "display_name", attrName+" edited"),
					resource.TestCheckResourceAttr(
						"split_traffic_type_attribute.foobar", "description", "this is my attribute description + edited"),
					resource.TestCheckResourceAttr(
						"split_traffic_type_attribute.foobar", "data_type", "NUMBER"),
					resource.TestCheckResourceAttr(
						"split_traffic_type_attribute.foobar", "suggested_values.#", "3"),
					resource.TestCheckResourceAttr(
						"split_traffic_type_attribute.foobar", "suggested_values.0", "1"),
					resource.TestCheckResourceAttr(
						"split_traffic_type_attribute.foobar", "suggested_values.1", "2"),
					resource.TestCheckResourceAttr(
						"split_traffic_type_attribute.foobar", "suggested_values.2", "3"),
					resource.TestCheckResourceAttr(
						"split_traffic_type_attribute.foobar", "is_searchable", "true"),
					resource.TestCheckResourceAttrSet(
						"split_traffic_type_attribute.foobar", "organization_id"),
				),
			},
		},
	})
}

func testAccCheckSplitTrafficTypeAttribute_basic(workspaceID, attrIdentifier, ttName, attrName string) string {
	return fmt.Sprintf(`
provider "split" {
	remove_environment_from_state_only = true
}

resource "split_traffic_type" "foobar" {
	workspace_id = "%[1]s"
	name = "%[3]s"
}

resource "split_traffic_type_attribute" "foobar" {
	workspace_id = "%[1]s"
	traffic_type_id = split_traffic_type.foobar.id
	identifier = "%[2]s"
	display_name = "%[4]s"
	description = "this is my attribute description"
	data_type = "STRING"
	suggested_values = ["a", "b", "c"]
	is_searchable = true
}
`, workspaceID, attrIdentifier, ttName, attrName)
}

func testAccCheckSplitTrafficTypeAttribute_updates(workspaceID, attrIdentifier, ttName, attrName string) string {
	return fmt.Sprintf(`
provider "split" {
	remove_environment_from_state_only = true
}

resource "split_traffic_type" "foobar" {
	workspace_id = "%[1]s"
	name = "%[3]s"
}

resource "split_traffic_type_attribute" "foobar" {
	workspace_id = "%[1]s"
	traffic_type_id = split_traffic_type.foobar.id
	identifier = "%[2]s"
	display_name = "%[4]s edited"
	description = "this is my attribute description + edited"
	data_type = "NUMBER"
	suggested_values = ["1", "2", "3"]
	is_searchable = true
}
`, workspaceID, attrIdentifier, ttName, attrName)
}
