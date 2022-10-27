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

### Static credentials

Credentials can be provided statically by adding an `api_key` arguments to the Split provider block:

```hcl
provider "split" {
  api_key = "SOME_API_KEY"
}
```

### Environment variables

When the Split provider block does not contain an `api_key` argument, the missing credentials will be sourced
from the environment via the `SPLIT_API_KEY` environment variables respectively:

```hcl
provider "split" {}
```

```shell
$ export SPLIT_API_KEY="SOME_KEY"
$ terraform plan
Refreshing Terraform state in-memory prior to plan...
```

## Rate Limiting

The Split provider provides automatic backoff in the event the provider detects the Split API has
[rate limited](https://docs.split.io/reference/rate-limiting) your Terraform operations. Please note
this backoff has a max timeout that can be configured by [`client_timeout`](#argument-reference) within
your `provider {}` block.

## Argument Reference

The following arguments are supported:

* `api_key` - (Required) Split API key. It must be provided, but it can also
  be sourced from [other locations](#Authentication).

* `base_url` - (Optional) Custom API URL.
  Can also be sourced from the `SPLIT_API_URL` environment variable.

* `remove_environment_from_state_only` - (Optional) Configure `split_environment` to only remove the resource from
  state upon deletion. This is to address out-of-band, UI based prerequisites Split has when deleting an environment.
  Defaults to `false`.

* `client_timeout` - (Optional) Configure client (http) timeout before aborting. This is to address the client retrying forever.
  It's expressed in an integer that represents seconds. Defaults to `300` seconds, or `5` minutes.