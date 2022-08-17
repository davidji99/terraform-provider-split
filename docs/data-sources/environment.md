---
layout: "split"
page_title: "Split: split_environment"
sidebar_current: "docs-split-datasource-environment"
description: |-
Get information about a Split environment
---

# Data Source: split_environment

Use this data source to get information about a Split [environment](https://help.split.io/hc/en-us/articles/360019915771-Environments).

## Example Usage

```hcl-terraform
data "split_environment" "production" {
  workspace_id = "71572aa0-3177-4591-946c-6bd4a7197cdb"
  name = "user"
}
```

## Argument Reference

The following arguments are supported:

* `workspace_id` - (Required) `<string>` The UUID of the workspace
* `name` - (Required) `<string>` Name of the environment

## Attributes Reference

The following attributes are exported:

* `production` - If the environment is production
