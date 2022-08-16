---
layout: "split"
page_title: "Split: split_environment"
sidebar_current: "docs-split-resource-environment"
description: |-
  Provides the ability to manage a Split environment.
---

# split_environment

This resource provides the ability to manage an [Environment](https://help.split.io/hc/en-us/articles/360019915771-Environments).
Environments allow you to manage your splits throughout your development lifecycle — from local development to staging and production.

## Example Usage

```hcl-terraform
data "split_workspace" "default" {
  name = "default"
}

resource "split_environment" "foobar" {
  workspace_id = data.split_workspace.default.id
  name = "production-canary"
  production = true
}
```

## Argument Reference

The following arguments are supported:

* `workspace_id` - (Required) `<string>` The UUID of the workspace you want to create the environment in.
* `name` - (Required) `<string>` Name of the environment.
* `production` - (Optional) `<boolean>` Whether the environment is deemed 'production'. Defaults to `false`.

## Attributes Reference

The following attributes are exported:

n/a

## Import

An existing environment can be imported using the combination of the workspace UUID
and environment ID separated by a colon (':').

For example:

```shell script
$ terraform import split_environment.foobar "0b46d8f7-9435-4f74-a770-3fcb22fbbfe6:110b3876-1d38-11ed-861d-0242ac120002"
```