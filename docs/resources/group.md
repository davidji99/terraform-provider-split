---
layout: "split"
page_title: "Split: split_group"
sidebar_current: "docs-split-resource-group"
description: |-
Provides the ability to manage a Split group.
---

# split_group

This resource provides the ability to manage a group in Split.

-> **IMPORTANT!**
Groups are not available on Split's free tier.

## Example Usage

```hcl-terraform
resource "split_group" "foobar" {
  name = "name_of_my_group"
  description = "description_of_my_group"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) `<string>` Name of the group.
* `description` - (Optional) `<boolean>` Description of the group.

## Attributes Reference

The following attributes are exported:

n/a

## Import

An existing group can be imported using the group UUID.

For example:

```shell script
$ terraform import split_group.foobar "0b46d8f7-9435-4f74-a770-3fcb22fbbfe6"
```