package split

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSplitSplit_Basic(t *testing.T) {
	workspaceID := testAccConfig.GetWorkspaceIDorSkip(t)
	trafficTypeName := fmt.Sprintf("tt-tftest-%s", acctest.RandString(10))
	splitName := fmt.Sprintf("s-tftest-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSplitSplit_basic(workspaceID, trafficTypeName, splitName, "my split description"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"split_split.foobar", "workspace_id", workspaceID),
					resource.TestCheckResourceAttr(
						"split_split.foobar", "name", splitName),
					resource.TestCheckResourceAttr(
						"split_split.foobar", "description", "my split description"),
					resource.TestCheckResourceAttrSet(
						"split_split.foobar", "traffic_type_id"),
				),
			},
			{
				Config: testAccCheckSplitSplit_basic(workspaceID, trafficTypeName, splitName, "my split edited description"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"split_split.foobar", "workspace_id", workspaceID),
					resource.TestCheckResourceAttr(
						"split_split.foobar", "name", splitName),
					resource.TestCheckResourceAttr(
						"split_split.foobar", "description", "my split edited description"),
					resource.TestCheckResourceAttrSet(
						"split_split.foobar", "traffic_type_id"),
				),
			},
		},
	})
}

func TestAccSplitSplit_WithOwners(t *testing.T) {
	workspaceID := testAccConfig.GetWorkspaceIDorSkip(t)
	ownerID := testAccConfig.GetUserIDorSkip(t)
	trafficTypeName := fmt.Sprintf("tt-tftest-%s", acctest.RandString(10))
	splitName := fmt.Sprintf("s-tftest-%s", acctest.RandString(10))
	t.Logf("Feature flag name: %s", splitName)
	t.Logf("Workspace ID: %s", workspaceID)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSplitSplit_withOwners(workspaceID, trafficTypeName, splitName, "split with owners", ownerID, "user"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("split_split.foobar", "workspace_id", workspaceID),
					resource.TestCheckResourceAttr("split_split.foobar", "name", splitName),
					resource.TestCheckResourceAttr("split_split.foobar", "description", "split with owners"),
					resource.TestCheckResourceAttrSet("split_split.foobar", "traffic_type_id"),
					resource.TestCheckResourceAttr("split_split.foobar", "owners.#", "1"),
					resource.TestCheckResourceAttr("split_split.foobar", "owners.0.id", ownerID),
					resource.TestCheckResourceAttr("split_split.foobar", "owners.0.type", "user"),
				),
			},
		},
	})
}

func TestAccSplitSplit_WithGroupOwner(t *testing.T) {
	workspaceID := testAccConfig.GetWorkspaceIDorSkip(t)
	groupID := testAccConfig.GetGroupIDorSkip(t)
	trafficTypeName := fmt.Sprintf("tt-tftest-%s", acctest.RandString(10))
	splitName := fmt.Sprintf("s-tftest-%s", acctest.RandString(10))
	t.Logf("Feature flag name: %s", splitName)
	t.Logf("Workspace ID: %s", workspaceID)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSplitSplit_withOwners(workspaceID, trafficTypeName, splitName, "split with group owner", groupID, "group"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("split_split.foobar", "workspace_id", workspaceID),
					resource.TestCheckResourceAttr("split_split.foobar", "name", splitName),
					resource.TestCheckResourceAttr("split_split.foobar", "description", "split with group owner"),
					resource.TestCheckResourceAttrSet("split_split.foobar", "traffic_type_id"),
					resource.TestCheckResourceAttr("split_split.foobar", "owners.#", "1"),
					resource.TestCheckResourceAttr("split_split.foobar", "owners.0.id", groupID),
					resource.TestCheckResourceAttr("split_split.foobar", "owners.0.type", "group"),
				),
			},
		},
	})
}

func TestAccSplitSplit_InvalidName(t *testing.T) {
	workspaceID := testAccConfig.GetWorkspaceIDorSkip(t)
	trafficTypeName := fmt.Sprintf("tt-tftest-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckSplitSplit_basic(workspaceID, trafficTypeName, "1invalidname", "my split description"),
				ExpectError: regexp.MustCompile(`Name must start with a letter and can contain hyphens, underscores, letters, and numbers`),
			},
		},
	})
}

func testAccCheckSplitSplit_basic(workspaceID, trafficTypeName, splitName, splitDescription string) string {
	return fmt.Sprintf(`
provider "split" {
	remove_environment_from_state_only = true
}

resource "split_traffic_type" "foobar" {
	workspace_id = "%[1]s"
	name = "%[2]s"
}

resource "split_split" "foobar" {
	workspace_id = "%[1]s"
	traffic_type_id = split_traffic_type.foobar.id
	name = "%[3]s"
	description = "%[4]s"
}
`, workspaceID, trafficTypeName, splitName, splitDescription)
}

// testAccCheckSplitSplit_withOwners builds config with one owner. ownerType is "user" or "group".
func testAccCheckSplitSplit_withOwners(workspaceID, trafficTypeName, splitName, splitDescription, ownerID, ownerType string) string {
	return fmt.Sprintf(`
provider "split" {
	remove_environment_from_state_only = true
}

resource "split_traffic_type" "foobar" {
	workspace_id = "%[1]s"
	name = "%[2]s"
}

resource "split_split" "foobar" {
	workspace_id   = "%[1]s"
	traffic_type_id = split_traffic_type.foobar.id
	name          = "%[3]s"
	description   = "%[4]s"
	owners {
		id   = "%[5]s"
		type = "%[6]s"
	}
}
`, workspaceID, trafficTypeName, splitName, splitDescription, ownerID, ownerType)
}
