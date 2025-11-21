---
layout: "split"
page_title: "Split: split_flag_set"
sidebar_current: "docs-split-datasource-flag-set"
description: |-
Get information about a Split flag set
---

# Data Source: split_flag_set

Use this data source to get information about a Split [flag set](https://help.split.io/hc/en-us/articles/360019916691-Flag-sets).

## Example Usage

```hcl-terraform
data "split_flag_set" "mobile_features" {
  workspace_id = "71572aa0-3177-4591-946c-6bd4a7197cdb"
  name = "mobile-features"
}
```

## Argument Reference

The following arguments are supported:

* `workspace_id` - (Required) `<string>` The UUID of the workspace
* `name` - (Required) `<string>` Name of the flag set

## Attributes Reference

The following attributes are exported:

* `description` - Description of the flag set
