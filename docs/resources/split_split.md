---
layout: "split"
page_title: "Split: split_split"
sidebar_current: "docs-split-resource-split"
description: |-
Provides the ability to manage a split.
---

# split_split

This resource provides the ability to manage a Split. A Split is a feature flag, toggle, or experiment.

## Example Usage

```hcl-terraform
data "split_workspace" "default" {
  name = "default"
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
```

## Argument Reference

The following arguments are supported:

* `workspace_id` - (Required) `<string>` The UUID of the workspace.
* `traffic_type_id` - (Required) `<string>` The UUID of the traffic type.
* `name` - (Required) `<string>` Name of Split. Name must start with a letter and can contain hyphens, underscores, letters, and numbers
* `description` - (Optional) `<string>` Description of Split.

## Attributes Reference

The following attributes are exported:

n/a

## Import

An existing split can be imported using the combination of the workspace UUID
and split UUID separated by a colon (':').

For example:

```shell script
$ terraform import split_split.foobar "0b46d8f7-9435-4f74-a770-3fcb22fbbfe6:8e52ce80-e05b-11ec-800d-5a826ff9ecd9"
```