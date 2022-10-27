package split

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccSplitSplitDefinition_Basic(t *testing.T) {
	workspaceID := testAccConfig.GetWorkspaceIDorSkip(t)
	envID := testAccConfig.GetEnvironmentIDorSkip(t)
	trafficTypeID := testAccConfig.GetTrafficTypeIDorSkip(t)
	trafficTypeName := fmt.Sprintf("tt-tftest-%s", acctest.RandString(10))
	splitName := fmt.Sprintf("s-tftest-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSplitSplitDefinition_basic(workspaceID, trafficTypeName, splitName,
					"my split description", envID, trafficTypeID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"split_split_definition.foobar", "workspace_id", workspaceID),
					resource.TestCheckResourceAttr(
						"split_split_definition.foobar", "split_name", splitName),
					resource.TestCheckResourceAttr(
						"split_split_definition.foobar", "environment_id", envID),
					resource.TestCheckResourceAttr(
						"split_split_definition.foobar", "default_treatment", "treatment_123"),
					resource.TestCheckResourceAttr(
						"split_split_definition.foobar", "treatment.0.name", "treatment_123"),
					resource.TestCheckResourceAttr(
						"split_split_definition.foobar", "treatment.0.configurations", "{\"key\":\"value\"}"),
					resource.TestCheckResourceAttr(
						"split_split_definition.foobar", "treatment.0.description", "my treatment 123"),
					resource.TestCheckResourceAttr(
						"split_split_definition.foobar", "treatment.1.name", "treatment_456"),
					resource.TestCheckResourceAttr(
						"split_split_definition.foobar", "treatment.1.configurations", "{\"key2\":\"value2\"}"),
					resource.TestCheckResourceAttr(
						"split_split_definition.foobar", "treatment.1.description", "my treatment 456"),
					resource.TestCheckResourceAttr(
						"split_split_definition.foobar", "default_rule.0.treatment", "treatment_123"),
					resource.TestCheckResourceAttr(
						"split_split_definition.foobar", "default_rule.0.size", "100"),
					resource.TestCheckResourceAttr(
						"split_split_definition.foobar", "rule.0.bucket.0.treatment", "treatment_123"),
					resource.TestCheckResourceAttr(
						"split_split_definition.foobar", "rule.0.bucket.0.size", "100"),
					resource.TestCheckResourceAttr(
						"split_split_definition.foobar", "rule.0.condition.0.combiner", "AND"),
					resource.TestCheckResourceAttr(
						"split_split_definition.foobar", "rule.0.condition.0.matcher.0.type", "EQUAL_SET"),
					resource.TestCheckResourceAttr(
						"split_split_definition.foobar", "rule.0.condition.0.matcher.0.attribute", "test_string"),
					resource.TestCheckResourceAttrSet(
						"split_split_definition.foobar", "rule.0.condition.0.matcher.0.strings.#"),
				),
			},
			{
				Config: testAccCheckSplitSplitDefinition_basicUpdated(workspaceID, trafficTypeName, splitName,
					"my split description", envID, trafficTypeID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"split_split_definition.foobar", "workspace_id", workspaceID),
					resource.TestCheckResourceAttr(
						"split_split_definition.foobar", "split_name", splitName),
					resource.TestCheckResourceAttr(
						"split_split_definition.foobar", "environment_id", envID),
					resource.TestCheckResourceAttr(
						"split_split_definition.foobar", "default_treatment", "treatment_123"),
					resource.TestCheckResourceAttr(
						"split_split_definition.foobar", "treatment.0.name", "treatment_123"),
					resource.TestCheckResourceAttr(
						"split_split_definition.foobar", "treatment.0.configurations", "{\"key\":\"value\"}"),
					resource.TestCheckResourceAttr(
						"split_split_definition.foobar", "treatment.0.description", "my treatment 123"),
					resource.TestCheckResourceAttr(
						"split_split_definition.foobar", "treatment.1.name", "treatment_789"),
					resource.TestCheckResourceAttr(
						"split_split_definition.foobar", "treatment.1.configurations", "{\"key3\":\"value3\"}"),
					resource.TestCheckResourceAttr(
						"split_split_definition.foobar", "treatment.1.description", "my treatment 789"),
					resource.TestCheckResourceAttr(
						"split_split_definition.foobar", "default_rule.0.treatment", "treatment_789"),
					resource.TestCheckResourceAttr(
						"split_split_definition.foobar", "default_rule.0.size", "100"),
					resource.TestCheckResourceAttr(
						"split_split_definition.foobar", "rule.0.bucket.0.treatment", "treatment_789"),
					resource.TestCheckResourceAttr(
						"split_split_definition.foobar", "rule.0.bucket.0.size", "100"),
					resource.TestCheckResourceAttr(
						"split_split_definition.foobar", "rule.0.condition.0.combiner", "AND"),
					resource.TestCheckResourceAttr(
						"split_split_definition.foobar", "rule.0.condition.0.matcher.0.type", "EQUAL_SET"),
					resource.TestCheckResourceAttr(
						"split_split_definition.foobar", "rule.0.condition.0.matcher.0.attribute", "test_string"),
					resource.TestCheckResourceAttrSet(
						"split_split_definition.foobar", "rule.0.condition.0.matcher.0.strings.#"),
				),
			},
		},
	})
}

func testAccCheckSplitSplitDefinition_basic(workspaceID, trafficTypeName, splitName, splitDescription, envID, trafficTypeID string) string {
	return fmt.Sprintf(`
provider "split" {
	remove_environment_from_state_only = true
}

resource "split_split" "foobar" {
	workspace_id = "%[1]s"
	traffic_type_id = "%[6]s"
	name = "%[3]s"
	description = "%[4]s"
}

resource "split_split_definition" "foobar" {
	workspace_id = "%[1]s"
	split_name = split_split.foobar.name
	environment_id = "%[5]s"

	default_treatment = "treatment_123"
	treatment {
		name = "treatment_123"
		configurations = "{\"key\":\"value\"}"
		description = "my treatment 123"
	}
	treatment {
		name = "treatment_456"
		configurations = "{\"key2\":\"value2\"}"
		description = "my treatment 456"
	}

	default_rule {
		treatment = "treatment_123"
		size = 60
	}

	default_rule {
		treatment = "treatment_123"
		size = 40
	}

	rule {
		bucket {
			treatment = "treatment_123"
			size = 100
		}
		condition {
			combiner = "AND"
			matcher {
				type = "EQUAL_SET"
				attribute = "test_string"
				strings = ["test_string"]
			}
		}
	}
}
`, workspaceID, trafficTypeName, splitName, splitDescription, envID, trafficTypeID)
}

func testAccCheckSplitSplitDefinition_basicUpdated(workspaceID, trafficTypeName, splitName, splitDescription, envID, trafficTypeID string) string {
	return fmt.Sprintf(`
provider "split" {
	remove_environment_from_state_only = true
}

resource "split_split" "foobar" {
	workspace_id = "%[1]s"
	traffic_type_id = "%[6]s"
	name = "%[3]s"
	description = "%[4]s"
}

resource "split_split_definition" "foobar" {
	workspace_id = "%[1]s"
	split_name = split_split.foobar.name
	environment_id = "%[5]s"

	default_treatment = "treatment_123"
	treatment {
		name = "treatment_123"
		configurations = "{\"key\":\"value\"}"
		description = "my treatment 123"
	}
	treatment {
		name = "treatment_789"
		configurations = "{\"key3\":\"value3\"}"
		description = "my treatment 789"
	}

	default_rule {
		treatment = "treatment_789"
		size = 100
	}
	rule {
		bucket {
			treatment = "treatment_789"
			size = 100
		}
		condition {
			combiner = "AND"
			matcher {
				type = "EQUAL_SET"
				attribute = "test_string"
				strings = ["test_string"]
			}
		}
	}
}
`, workspaceID, trafficTypeName, splitName, splitDescription, envID, trafficTypeID)
}
