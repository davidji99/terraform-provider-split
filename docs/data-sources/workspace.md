---
layout: "split"
page_title: "Split: split_workspace"
sidebar_current: "docs-split-datasource-workspace"
description: |-
    Get information about a Split workspace
---

# Data Source: split_workspace

Use this data source to get information about a Split workspace.

## Example Usage

```hcl-terraform
data "split_environment" "default" {
  name = "default"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) `<string>` Name of the workspace

## Attributes Reference

The following attributes are exported:

* `type` - Type of workspace.

* `requires_title_and_comments` - Whether the workspace requires titles and comments.