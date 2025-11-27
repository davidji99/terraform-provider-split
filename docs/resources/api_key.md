---
layout: "split"
page_title: "Split: split_api_key"
sidebar_current: "docs-split-resource-api-key"
description: |-
  Provides the ability to manage a Split API key.
---

# split_api_key

This resource provides the ability to manage an [API key](https://docs.split.io/reference/api-keys-overview).

Due to API limitations, it is not possible update an existing API key. Any modifications to a `split_api_key` resource
will result in a destroy and recreate process.

-> **DEPRECATION NOTICE**
When using `harness_token` for authentication (x-api-key header), API keys with `type = "admin"` are deprecated and cannot be used. Please use the Harness Terraform provider instead for managing admin API keys when using Harness authentication.

-> **IMPORTANT!**
Please be very careful when deleting this resource as the deleted API keys are NOT recoverable and invalidated immediately.
Furthermore, this resource renders the actual API key plain-text in your state file.
Please ensure that your state file is properly secured and encrypted at rest.

## Example Usage

```hcl-terraform
data "split_workspace" "default" {
  name = "default"
}

resource "split_environment" "foobar" {
	workspace_id = data.split_workspace.default.id
	name = "staging"
	production = true
}

resource "split_api_key" "foobar" {
	workspace_id = data.split_workspace.default.id
	name = "my client side key"
	type = "client_side"
	environment_ids = [split_environment.foobar.id]
}
```

## Argument Reference

The following arguments are supported:

* `workspace_id` - (Required) `<string>` The UUID of the workspace you want to create the environment in.
* `environment_ids` - (Required) `<list(string)>` List of environment UUIDs.
* `name` - (Required) `<boolean>` Name of the API key.
* `type` - (Required) `<boolean>` Type of the API key. Refer to Split [documentation](https://docs.split.io/reference/create-an-api-key#supported-types) on acceptable values, case sensitive.
* `roles` - (Required) `<boolean>` Supported only when `type=admin` API keys. For the full list of allowed Admin API Key roles, refer
to Split [documentation](https://docs.split.io/reference/api-keys-overview#admin-api-key-roles), case sensitive.

## Attributes Reference

The following attributes are exported:

n/a

## Import

Due to Split API limitations, it is not possible to import an existing API key.
