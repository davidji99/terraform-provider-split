---
layout: "split"
page_title: "Split: split_flag_set"
sidebar_current: "docs-split-resource-flag-set"
description: |-
  Provides the ability to manage a Split flag set.
---

# split_flag_set

This resource provides the ability to manage a [Flag Set](https://help.split.io/hc/en-us/articles/360019916691-Flag-sets).
Flag sets allow you to group feature flags together for easier management and organization.

## Example Usage

```hcl-terraform
data "split_workspace" "default" {
  name = "default"
}

resource "split_flag_set" "example" {
  workspace_id = data.split_workspace.default.id
  name = "mobile-features"
  description = "Feature flags for mobile applications"
}
```

## Argument Reference

The following arguments are supported:

* `workspace_id` - (Required) `<string>` The UUID of the workspace you want to create the flag set in.
* `name` - (Required) `<string>` Name of the flag set. Must be between 1 and 50 characters long, start with a lowercase letter or digit, and may contain only lowercase letters, digits, underscores, and dots.
* `description` - (Optional) `<string>` Description of the flag set.

## Attributes Reference

The following attributes are exported:

n/a

## Import

An existing flag set can be imported using its UUID.

For example:

```shell script
$ terraform import split_flag_set.example "110b3876-1d38-11ed-861d-0242ac120002"
```
