package split

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccDataSourceSplitFlagSet_Basic(t *testing.T) {
	workspaceID := testAccConfig.GetWorkspaceIDorSkip(t)
	name := fmt.Sprintf("tftest-%s", acctest.RandString(8))
	description := fmt.Sprintf("Test description for %s", name)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSplitFlagSet_basic(workspaceID, name, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.split_flag_set.test", "workspace_id", workspaceID),
					resource.TestCheckResourceAttr("data.split_flag_set.test", "name", name),
					resource.TestCheckResourceAttr("data.split_flag_set.test", "description", description),
					resource.TestCheckResourceAttrSet("data.split_flag_set.test", "id"),
				),
			},
		},
	})
}

func TestAccDataSourceSplitFlagSet_BasicNoDescription(t *testing.T) {
	workspaceID := testAccConfig.GetWorkspaceIDorSkip(t)
	name := fmt.Sprintf("tftest-%s", acctest.RandString(8))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSplitFlagSet_basicNoDescription(workspaceID, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.split_flag_set.test", "workspace_id", workspaceID),
					resource.TestCheckResourceAttr("data.split_flag_set.test", "name", name),
					resource.TestCheckResourceAttrSet("data.split_flag_set.test", "id"),
				),
			},
		},
	})
}

func testAccDataSourceSplitFlagSet_basic(workspaceID, name, description string) string {
	return fmt.Sprintf(`
resource "split_flag_set" "foobar" {
	workspace_id = "%s"
	name = "%s"
	description = "%s"
}

data "split_flag_set" "test" {
	workspace_id = split_flag_set.foobar.workspace_id
	name = split_flag_set.foobar.name
}
`, workspaceID, name, description)
}

func testAccDataSourceSplitFlagSet_basicNoDescription(workspaceID, name string) string {
	return fmt.Sprintf(`
resource "split_flag_set" "foobar" {
	workspace_id = "%s"
	name = "%s"
}

data "split_flag_set" "test" {
	workspace_id = split_flag_set.foobar.workspace_id
	name = split_flag_set.foobar.name
}
`, workspaceID, name)
}
