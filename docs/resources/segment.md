---
layout: "split"
page_title: "Split: split_segment"
sidebar_current: "docs-split-resource-segment"
description: |-
Provides the ability to manage a Split segment.
---

# split_segment

This resource provides the ability to manage a [Segment](https://help.split.io/hc/en-us/articles/360020407512-Create-a-segment).
A segment is a pre-defined group of customers that a feature can be targeted to. Segments are best for targeting relatively
fixed or specific groups of users that you can easily identify, like a whitelist of accounts.

## Example Usage

```hcl-terraform
data "split_environment" "default" {
  name = "default"
}

data "split_traffic_type" "user" {
  workspace_id = data.split_environment.default.id
  name = "user"
}

resource "split_segment" "foobar" {
  workspace_id = data.split_environment.default.id
  traffic_type_id = data.split_traffic_type.user.id
  name = "name_of_my_segment"
  description = "description_of_my_segment"
}
```

## Argument Reference

The following arguments are supported:

* `workspace_id` - (Required) `<string>` The UUID of the workspace.
* `traffic_type_id` - (Required) `<string>` The UUID of the traffice type.
* `name` - (Required) `<string>` Name of the segment.
* `description` - (Optional) `<boolean>` Description of the segment.

## Attributes Reference

The following attributes are exported:

n/a

## Import

An existing segment can be imported using the combination of the workspace UUID
and segment name separated by a colon (':').

For example:

```shell script
$ terraform import split_segment.foobar "0b46d8f7-9435-4f74-a770-3fcb22fbbfe6:name_of_my_segment"
```