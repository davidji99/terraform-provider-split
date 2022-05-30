---
layout: "split"
page_title: "Split: split_traffic_type"
sidebar_current: "docs-split-resource-traffic-type"
description: |-
Provides the ability to manage a Split traffic type.
---

# split_traffic_type

This resource provides the ability to manage a traffic type.

A traffic type is a particular identifier type for any hierarchy of your customer base. Traffic types in Split are
completely customizable and can be any database key you choose to send to Split, i.e. a user ID, account ID, IP address,
browser ID, etc. Essentially, any internal database key you're using to track what "customer" means to you.

## Example Usage

```hcl-terraform
data "split_workspace" "default" {
  name = "default"
}

resource "split_traffic_type" "foobar" {
  workspace_id = data.split_workspace.default.id
  name = "name_of_my_traffic_type"
}
```

## Argument Reference

The following arguments are supported:

* `workspace_id` - (Required) `<string>` The UUID of the workspace.
* `name` - (Required) `<string>` Name of the traffic type.

## Attributes Reference

The following attributes are exported:

* `type`  - Type of traffic type.
* `display_attribute_id` - Attribute used for display name in UI.

## Import

An existing traffic type can be imported using the combination of the workspace UUID
and traffic type UUID separated by a colon (':').

For example:

```shell script
$ terraform import split_traffic_type.foobar "0b46d8f7-9435-4f74-a770-3fcb22fbbfe6:268103c0-8eef-11ec-a34a-1ae28fca10c9"
```