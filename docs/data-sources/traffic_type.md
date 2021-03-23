---
layout: "split"
page_title: "Split: split_traffic_type"
sidebar_current: "docs-split-datasource-traffic-type"
description: |-
Get information about a Split traffic type
---

# Data Source: split_traffic_type

Use this data source to get information about a Split [traffic type](https://help.split.io/hc/en-us/articles/360019916311-Traffic-type#:~:text=Split%20allows%20you%20to%20have,needed%20during%20your%20account%20setup.).

## Example Usage

```hcl-terraform
data "split_traffic_type" "user" {
  workspace_id = "71572aa0-3177-4591-946c-6bd4a7197cdb"
  name = "user"
}
```

## Argument Reference

The following arguments are supported:

* `workspace_id` - (Required) `<string>` The UUID of the workspace.
* `name` - (Required) `<string>` Name of the traffic type.

## Attributes Reference

The following attributes are exported:

* `type` - Type of traffic type.
* `display_attribute_id` - The traffic type display attribute ID.
* `workspace_name` - Name of the workspace.