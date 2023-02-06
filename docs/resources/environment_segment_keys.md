---
layout: "split"
page_title: "Split: split_environment"
sidebar_current: "docs-split-resource-environment"
description: |-
  Provides the ability to manage a Segment keys in a Split environment.
---

# split_environment

This resource pProvides the ability to manage a Segment keys in a Split environment.

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

resource "split_environment_segment_keys" "foobar" {
	environment_id = split_environment.foobar.id
	segment_name = split_segment_environment_association.foobar.segment_name
	keys = ["mytestkey1", "mytestkey2"]
}
```

## Argument Reference

The following arguments are supported:

* `environment_id` - (Required) `<string>` The UUID of the environment.
* `segment_name` - (Required) `<string>` Name of the segment.
* `keys` - (Optional) `<list(string)>` List of identifiers, aka keys. Can only add up to 10,000 keys at a time.
  Order of keys does not matter.

## Attributes Reference

The following attributes are exported:

n/a

## Import

An existing environment can be imported using the combination of the environment UUID
and segment name separated by a colon (':').

For example:

```shell script
$ terraform import split_environment_segment_keys.foobar "0b46d8f7-9435-4f74-a770-3fcb22fbbfe6:name_of_my_segment"
```
