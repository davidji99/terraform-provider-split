---
layout: "split"
page_title: "Split: split_traffic_type_attribute"
sidebar_current: "docs-split-resource-traffic-type-attribute"
description: |-
Provides the ability to manage a Split traffic type attribute.
---

# split_traffic_type_attribute

This resource provides the ability to manage an [Traffic type Attribute](https://help.split.io/hc/en-us/articles/360020793231-Target-with-custom-attributes).

Attribute schemas (commonly called Attributes) define are definitions for attribute values, which are stored on Traffic
ID Keys (Identities). With the API you can attributes programmatically manage these attributes.

## Default attributes

Attributes are also created automatically when new keys are saved to Split. Those newly created attributes are
created with an unknown data type and with no display or description information. You will need to import these
attributes first before managing them via Terraform.

## Example Usage

```hcl-terraform
resource "split_traffic_type" "foobar" {
	workspace_id = "my_workspace_id"
	name = "my_workspace_name"
}

resource "split_traffic_type_attribute" "foobar" {
	workspace_id = "my_workspace_id"
	traffic_type_id = split_traffic_type.foobar.id
	identifier = "my-attribute"
	display_name = "my attribute display name"
	description = "this is my attribute description"
	data_type = "NUMBER"
	suggested_values = ["1", "2", "3"]
	is_searchable = true
}
```

## Argument Reference

The following arguments are supported:

* `workspace_id` - (Required) `<string>` The UUID of the workspace.
* `traffic_type_id` - (Required) `<string>` The UUID of the traffic type.
* `identifier` - (Required) `<string>` Id of the attribute. Max length is 200 characters.
* `display_name` - (Required) `<string>` Display name of the attribute. Max length is 200 characters.
* `description` - (Required) `<string>` Description of the attribute. Max length is 500 characters.
* `data_type` - (Optional) `<string>` Data type of the attribute. See the [specification](#data_type) below for more details.
* `suggested_values` - (Optional) `<list(string)>` Suggested value(s) of the attribute. See the [specification](#suggested_values) below for more details.
* `is_searchable` - (Optional) `<boolean>` Whether or not the attribute is searchable.

### `data_type`

Valid options are `STRING`, `DATETIME`, `NUMBER`, `SET` and are case sensitive. If you don't want
to explicitly set a type, do not define it in your resource configuration.

The data type is an optional parameter which is used for display formatting purposes.
This parameter is not used to validate the values on inbound keys, but incorrectly typed data will
be displayed as their raw string.

### `suggested_values`

Regardless of what [`data_type`](#data_type) is set to, all values need to be defined as strings.

The maximum number of values is 50 and each value cannot exceed 50 characters.

Values for datetime fields are expected to be passed as milliseconds past the epoch. Values for set fields are expected
to be passed in a string as comma separated values. If an unsupported data type is sent you will receive a 400 error code
as a response.

## Attributes Reference

The following attributes are exported:

* `organization_id` - The attribute's organization id.

## Import

An existing attribute can be imported using the combination of the workspace UUID, traffic type UUID,
and attribute ID separated by a colon (':').

For example:

```shell script
$ terraform import split_traffic_type_attribute.foobar "0b46d8f7-9435-4f74-a770-3fcb22fbbfe6:110b3876-1d38-11ed-861d-0242ac120002:my-attribute"
```
