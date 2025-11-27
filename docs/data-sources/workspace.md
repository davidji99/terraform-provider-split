---
layout: "split"
page_title: "Split: split_workspace"
sidebar_current: "docs-split-datasource-workspace"
description: |-
    Get information about a Split workspace
---

# Data Source: split_workspace

Use this data source to get information about a Split workspace.

-> **NOTE**
This data source is available when using either `api_key` or `harness_token` authentication. While the `split_workspace` resource is deprecated when using `harness_token`, this data source remains fully functional.

## Example Usage

```hcl-terraform
data "split_workspace" "default" {
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