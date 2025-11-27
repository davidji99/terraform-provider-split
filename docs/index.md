---
layout: "split"
page_title: "Provider: Split"
sidebar_current: "docs-split-index"
description: |-
  The Split provider is used to interact with resources provided by the Split API.
---

# Split Provider

The Split provider is used to interact with resources provided by the
[Split API](https://docs.split.io/reference#introduction).

## Contributing

Development happens in the [GitHub repo](https://github.com/davidji99/terraform-provider-split):

* [Releases](https://github.com/davidji99/terraform-provider-split/releases)
* [Issues](https://github.com/davidji99/terraform-provider-split/issues)

## Example Usage

```hcl
provider "split" {
}

# Create a new Split environment
resource "split_environment" "foobar" {
  # ...
}
```

## Authentication

The Split provider offers a flexible means of providing credentials for authentication.
The following methods are supported, listed in order of precedence, and explained below:

- Static credentials
- Environment variables

### Authentication Options

This provider supports two authentication methods:

1. **API Key Authentication (Default)**: Uses a Bearer token in the Authorization header.
   - Set via the `api_key` parameter or the `SPLIT_API_KEY` environment variable.

2. **Harness Token Authentication**: Uses the `x-api-key` header for authentication.
   - Set via the `harness_token` parameter or the `HARNESS_TOKEN` environment variable.
   - When this authentication method is used, the following resources are deprecated and cannot be used:
     - `split_user`
     - `split_group`
     - `split_workspace` (resource only - the `split_workspace` data source is still available)
     - `split_api_key` (only when `type = "admin"`)

### Static credentials

Credentials can be provided statically by adding an `api_key` or `harness_token` argument to the Split provider block:

```hcl
provider "split" {
  # Use either api_key (default) for Bearer token authentication
  api_key = "SOME_API_KEY"

  # OR use harness_token for x-api-key header authentication
  # harness_token = "SOME_HARNESS_TOKEN"
}
```

### Environment variables

When the Split provider block does not contain an `api_key` or `harness_token` argument, the missing credentials will be sourced
from the environment via the `SPLIT_API_KEY` or `HARNESS_TOKEN` environment variables respectively:

```hcl
provider "split" {}
```

```shell
# For API key authentication
$ export SPLIT_API_KEY="SOME_KEY"
$ terraform plan

# OR for Harness token authentication
$ export HARNESS_TOKEN="SOME_TOKEN"
$ terraform plan
```

## Rate Limiting

The Split provider provides automatic backoff in the event the provider detects the Split API has
[rate limited](https://docs.split.io/reference/rate-limiting) your Terraform operations. Please note
this backoff has a max timeout that can be configured by [`client_timeout`](#argument-reference) within
your `provider {}` block.

## Argument Reference

The following arguments are supported:

* `api_key` - (Optional) Split API key for Bearer token authentication. It can be provided here or
  sourced from the `SPLIT_API_KEY` environment variable. Either `api_key` or `harness_token` must be provided.

* `harness_token` - (Optional) Harness token for x-api-key header authentication. It can be provided here or
  sourced from the `HARNESS_TOKEN` environment variable. Either `api_key` or `harness_token` must be provided.
  When using `harness_token`, certain resources are deprecated (see [Authentication Options](#authentication-options)).

* `base_url` - (Optional) Custom API URL.
  Can also be sourced from the `SPLIT_API_URL` environment variable.

* `remove_environment_from_state_only` - (Optional) Configure `split_environment` to only remove the resource from
  state upon deletion. This is to address out-of-band, UI based prerequisites Split has when deleting an environment.
  Defaults to `false`.

* `client_timeout` - (Optional) Configure client (http) timeout before aborting. This is to address the client retrying forever.
  It's expressed in an integer that represents seconds. Defaults to `300` seconds, or `5` minutes.