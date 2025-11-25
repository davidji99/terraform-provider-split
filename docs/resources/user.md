---
layout: "split"
page_title: "Split: split_user"
sidebar_current: "docs-split-resource-user"
description: |-
Provides the ability to manage a Split user.
---

# split_user

This resource provides the ability to manage a user in Split.

-> **DEPRECATION NOTICE**
When using `harness_token` for authentication (x-api-key header), this resource is deprecated and cannot be used. Please use the Harness Terraform provider instead when using Harness authentication.

Due to API behavior, this resource does not provide the ability
to set a user's `name` attribute as you cannot update a user
that hasn't accepted its invitation. The resource cannot accurately
determine the user as accepted the invitation within a reasonable
amount of time.

### Deletion Behavior

Upon resource deletion, the following may occur:

* If the user has not accepted the invitation, the invitation will be deleted.
* If the user has accepted the invitation already, the user will be deactivated.

## Example Usage

```hcl-terraform
resource "split_user" "user" {
  email = "user@company.com"
}
```

## Argument Reference

The following arguments are supported:

* `email` - (Required) `<string>` Name of the user.

## Attributes Reference

The following attributes are exported:

* `name` - Name of the user. By default, the `name` is everything before the '@' in the user `email`.
* `2fa` - Whether 2FA is enabled.
* `status` - Status of the user.

## Import

An existing user can be imported using the user UUID.

For example:

```shell script
$ terraform import split_user.foobar "0b46d8f7-9435-4f74-a770-3fcb22fbbfe6"
```