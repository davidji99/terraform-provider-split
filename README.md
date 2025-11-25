Terraform Provider Split
=========================

This provider is used to configure certain resources supported by [Split API](https://docs.split.io/reference#introduction).

For provider bugs/questions, please open an issue on this repository.

Documentation
------------

Documentation about resources and data sources can be found
[here](https://registry.terraform.io/providers/davidji99/split/latest/docs).

Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) `v0.13.x`+
- [Go](https://golang.org/doc/install) 1.18 (to build the provider plugin)

Usage
-----

```hcl
provider "split" {
  version = "~> 0.1.0"
  
  # Use either api_key (default) for Bearer token authentication
  api_key = "YOUR_API_KEY"
  
  # OR use harness_token for x-api-key header authentication
  # harness_token = "YOUR_HARNESS_TOKEN"
}
```

### Authentication Options

This provider supports two authentication methods:

1. **API Key Authentication (Default)**: Uses a Bearer token in the Authorization header.
   - Set via the `api_key` parameter or the `SPLIT_API_KEY` environment variable.

2. **Harness Token Authentication**: Uses the `x-api-key` header for authentication.
   - Set via the `harness_token` parameter or the `HARNESS_TOKEN` environment variable.
   - When this authentication method is used, the following resources are deprecated and cannot be used:
     - `split_user`
     - `split_group`
     - `split_workspace`
     - `split_api_key` (only when `type = "admin"`)


Releases
------------

Provider binaries can be found [here](https://github.com/davidji99/terraform-provider-split/releases).

Development
-----------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.12+ is *required*).

If you wish to bump the provider version, you can do so in the file `version/version.go`.

### Build the Provider

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```shell script
$ make build
...
$ $GOPATH/bin/terraform-provider-split
...
```

### Using the Provider

To use the dev provider with local Terraform, copy the freshly built plugin into Terraform's local plugins directory:

```sh
cp $GOPATH/bin/terraform-provider-split ~/.terraform.d/plugins/
```

Set the split provider without a version constraint:

```hcl
provider "split" {}
```

Then, initialize Terraform:

```shell script
terraform init
```

### Testing

Please see the [TESTING](TESTING.md) guide for detailed instructions on running tests.

### Updating or adding dependencies

This project uses [Go Modules](https://github.com/golang/go/wiki/Modules) for dependency management.

This example will fetch a module at the release tag and record it in your project's `go.mod` and `go.sum` files.
It's a good idea to run `go mod tidy` afterward and then `go mod vendor` to copy the dependencies into a `vendor/` directory.

If a module does not have release tags, then `module@SHA` can be used instead.
