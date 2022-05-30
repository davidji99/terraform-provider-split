---
layout: "split"
page_title: "Split: split_workspace"
sidebar_current: "docs-split-resource-workspace"
description: |-
    Provides the ability to manage a Split workspace.
---

# split_workspace

Use this resource to manage a Split workspace.

Workspaces allow you to separately manage your feature flags and experiments across your different business units,
product lines, and applications.

## Example Usage

```hcl-terraform
resource "split_workspace" "foobar" {
  name = "my_new_workspace"
  require_title_comments = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) `<string>` Name of the workspace.
* `require_title_comments` - (Optional) `<boolean>` Require title and comments for splits, segment, and metric changes.

## Attributes Reference

The following attributes are exported:

n/a
