package split

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccSplitGroup_Basic(t *testing.T) {
	name := fmt.Sprintf("tftest-%s", acctest.RandString(8))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSplitGroup_basic(name, "created from Terraform"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"split_group.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"split_group.foobar", "description", "created from Terraform"),
				),
			},
			{
				Config: testAccCheckSplitGroup_basic(name+" edited", "created from Terraform edited"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"split_group.foobar", "name", name+" edited"),
					resource.TestCheckResourceAttr(
						"split_group.foobar", "description", "created from Terraform edited"),
				),
			},
		},
	})
}

func testAccCheckSplitGroup_basic(name, description string) string {
	return fmt.Sprintf(`
resource "split_group" "foobar" {
	name = "%s"
	description = "%s"
}
`, name, description)
}
