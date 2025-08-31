package split

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccSplitFlagSet_Basic(t *testing.T) {
	workspaceID := testAccConfig.GetWorkspaceIDorSkip(t)
	name := fmt.Sprintf("tftest-%s", acctest.RandString(8))
	description := fmt.Sprintf("Test description for %s", name)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSplitFlagSet_basic(workspaceID, name, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"split_flag_set.foobar", "workspace_id", workspaceID),
					resource.TestCheckResourceAttr(
						"split_flag_set.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"split_flag_set.foobar", "description", description),
				),
			},
		},
	})
}

func TestAccSplitFlagSet_BasicNoDescription(t *testing.T) {
	workspaceID := testAccConfig.GetWorkspaceIDorSkip(t)
	name := fmt.Sprintf("tftest-%s", acctest.RandString(8))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSplitFlagSet_basicNoDescription(workspaceID, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"split_flag_set.foobar", "workspace_id", workspaceID),
					resource.TestCheckResourceAttr(
						"split_flag_set.foobar", "name", name),
				),
			},
		},
	})
}

func testAccCheckSplitFlagSet_basic(workspaceID, name, description string) string {
	return fmt.Sprintf(`
resource "split_flag_set" "foobar" {
	workspace_id = "%s"
	name = "%s"
	description = "%s"
}
`, workspaceID, name, description)
}

func testAccCheckSplitFlagSet_basicNoDescription(workspaceID, name string) string {
	return fmt.Sprintf(`
resource "split_flag_set" "foobar" {
	workspace_id = "%s"
	name = "%s"
}
`, workspaceID, name)
}
