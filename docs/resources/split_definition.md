---
layout: "split"
page_title: "Split: split_split_definition"
sidebar_current: "docs-split-resource-split-definition"
description: |-
Provides the ability to manage a split definition.
---

# split_definition

This resource provides the ability to manage a Split Definition.

## Example Usage

```hcl-terraform
data "split_workspace" "default" {
  name = "default"
}

resource "split_environment" "foobar" {
  workspace_id = data.split_workspace.default.id
  name = "production-canary"
  production = true
}

resource "split_traffic_type" "foobar" {
  workspace_id = data.split_workspace.default.id
  name = "my_traffic_type"
}

resource "split_split" "foobar" {
  workspace_id = data.split_workspace.default.id
  traffic_type_id = split_traffic_type.foobar.id
  name = "my_split"
  description = "my split description"
}

resource "split_split_definition" "foobar" {
  workspace_id = data.split_workspace.default.id
  split_name = split_split.foobar.name
  environment_id = split_environment.foobar.id

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
    size = 100
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
```

## Argument Reference

The following arguments are supported:

* `workspace_id` - (Required) `<string>` The UUID of the workspace.
* `split_name` - (Required) `<string>` The name, not UUID, of the Split
* `environment_id` - (Required) `<string>` The UUID of the environment
* `default_treatment` - (Required) `<string>` Default treatment to place unassigned customers into or randomly distribute
  these customers between your treatments/variations based off of percentages you decide. This attribute value should
  match one of your `treatment.name` values.
* `treatment` - (Required) `<block>` See the [specification](#treatment) below for more details.
* `default_rule` - (Required) `<block>` See the [specification](#default_rule) below for more details.
* `rule` - (Required) `<block>` See the [specification](#rule) below for more details.

### `treatment`

What are the different variations of this feature? Define the treatments and add a description for each to quickly
explain the difference between each treatment. This attribute block supports the following:

* `name` - (Required) `<string>` Name of the treatment.
* `configurations` - (Required) `<string>` Dynamically configure components of your feature (e.g. A button's color or backend API pagination).
  This attribute's value must be a valid JSON string.
* `description` - (Optional) `<string>` Description of the treatment.
* `keys` - (Optional) `<list(string)>` List of target key ids.
* `segments` - (Optional) `<list(string)>` List of segments.

### `default_rule`

For any of your customers that weren't assigned a treatment in the sections above,
use this section to place them into treatments or randomly distribute these customers among your treatments/variations
based off of percentages you decide. This attribute block supports the following:

* `treatment` - (Required) `<string>` Name of a valid Treatment.
* `size` - (Required) `<integer>` Treatment size. The sum of all your `default_rule` blocks must equal `100`.

### `rule`

Target specific subsets of your customers based on a specific attribute (e.g., location or last-login date). These subsets
are placed in specific treatments or, based on the percentages you decide, we’ll take these subsets and randomly distribute
customers in the subset between your treatments. This attribute block supports the following:

* `bucket` - (Required) `<block>` List of treatments.
    * `treatment` - (Required) `<string>` Name of a valid Treatment.
    * `size` - (Required) `<integer>` Treatment size.
* `condition` - (Required) `<block>` Rule conditions.
    * `combiner` - (Required) `<string>` rule condition combiner.
    * `matcher` - (Required) `<block>` rule condition matcher.
        * `type` - (Required) `<string>` rule condition matcher type.
        * `attribute` - (Required) `<string>` rule condition matcher type.
        * `string` - (Optional) `<string>` This matcher selects customers with an attribute or key that matches the regex pattern set by this attribute.
        * `strings` - (Optional) `<list(string)>` rule condition matcher type.

It is recommended to view the UI in order to determine what are some of the possible attribute values.

## Attributes Reference

The following attributes are exported:

n/a

## Import

An existing split definition can be imported using the combination of the workspace UUID, split name, and environment UUID
separated by a colon (':').

For example:

```shell script
$ terraform import split_split_definition.foobar "0b46d8f7-9435-4f74-a770-3fcb22fbbfe6:my-split:8e52ce80-e05b-11ec-800d-5a826ff9ecd9"
```