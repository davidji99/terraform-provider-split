---
layout: "split"
page_title: "Split: split_attribute"
sidebar_current: "docs-split-resource-attribute"
description: |-
Provides the ability to manage a Split attribute.
---

# split_attribute

This resource provides the ability to manage an [Attribute](https://help.split.io/hc/en-us/articles/360020793231-Target-with-custom-attributes).

## Example Usage

```hcl-terraform
resource "split_traffic_type" "foobar" {
	workspace_id = "my_workspace_id"
	name = "my_workspace_name"
}

resource "split_attribute" "foobar" {
	workspace_id = "my_workspace_id"
	traffic_type_id = split_traffic_type.foobar.id
	identifier = "my-attribute"
	display_name = "my attribute display name"
	description = "this is my attribute description"
	data_type = "NUMBER"
	suggested_values = [1, 2, 3]
	is_searchable = true
}
```

## Argument Reference

The following arguments are supported:

* `workspace_id` - (Required) `<string>` The UUID of the workspace.
* `traffic_type_id` - (Required) `<string>` The UUID of the traffic type.
* `identifier` - (Required) `<string>` Id of the attribute.
* `display_name` - (Required) `<string>` Display name of the attribute.
* `description` - (Required) `<string>` Description of the attribute.
* `data_type` - (Required) `<string>` Data type of the attribute. Valid options are `STRING`, `DATETIME`, `NUMBER`, `SET`.
   Case sensitive.
* `suggested_values` - (Optional) `<list(string)>` Suggested value(s) of the attribute.
* `is_searchable` - (Optional) `<boolean>` Whether or not the attribute is searchable.

## Attributes Reference

The following attributes are exported:

* `organization_id` - The attribute's organization id.

## Import

An existing attribute can be imported using the combination of the workspace UUID, traffic type UUID,
and attribute ID separated by a colon (':').

For example:

```shell script
$ terraform import split_attribute.foobar "0b46d8f7-9435-4f74-a770-3fcb22fbbfe6:110b3876-1d38-11ed-861d-0242ac120002:my-attribute"
```
