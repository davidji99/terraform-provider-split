---
layout: "split"
page_title: "Split: split_segment_environment_association"
sidebar_current: "docs-split-resource-segment-environment-association"
description: |-
Provides the ability to enable and disable a Split segment in an environment.
---

# split_segment_environment_association

This resource provides the ability to enable and disable a Split segment in an environment.

## Example Usage

```hcl-terraform
data "split_workspace" "default" {
  name = "default"
}

resource "split_environment" "foobar" {
	workspace_id = data.split_workspace.default.id
	name = "foobar_env"
	production = true
}

resource "split_traffic_type" "foobar" {
	workspace_id = data.split_workspace.default.id
	name = "foobar_traffic_type_env"
}

resource "split_segment" "foobar" {
	workspace_id = data.split_workspace.default.id
	traffic_type_id = split_traffic_type.foobar.id
	name = "foobar segment"
	description = "description of my foobar segment"
}

resource "split_segment_environment_association" "foobar" {
	workspace_id = data.split_workspace.default.id
	environment_id = split_environment.foobar.id
	segment_name = split_segment.foobar.name
}
```

## Argument Reference

The following arguments are supported:

* `workspace_id` - (Required) `<string>` The UUID of the workspace.
* `environment_id` - (Required) `<string>` The UUID of the environment.
* `segment_name` - (Required) `<string>` Name of the segment.

## Attributes Reference

The following attributes are exported:

n/a

## Import

An existing segment can be imported using the combination of the workspace UUID, environment UUID, and
segment name separated by a colon (':').

For example:

```shell script
$ terraform import split_segment_environment_association.foobar "0b46d8f7-9435-4f74-a770-3fcb22fbbfe6:a6d5d991-4069-44bd-8d00-949d8cafd120:name_of_my_segment"
```
