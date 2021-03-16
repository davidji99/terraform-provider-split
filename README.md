Terraform Provider <REPLACE_ME>
=========================

Some general information about this provider.

Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) `v0.12.x`+
- [Go](https://golang.org/doc/install) 1.14 (to build the provider plugin)

Usage
-----

```hcl
provider "<REPLACE_ME>" {
  version = "~> 0.1.0"
}
```

This provider is not compatible with terraform `v0.11.x`.

Development
-----------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.14+ is *required*).

### Build the Provider

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```shell script
$ make build
...
$ $GOPATH/bin/terraform-provider-<REPLACE_ME>
...
```

### Using the Provider

To use the dev provider with local Terraform, copy the freshly built plugin into Terraform's local plugins directory:

```sh
cp $GOPATH/bin/terraform-provider-<REPLACE_ME> ~/.terraform.d/plugins/
```

Set the <REPLACE_ME> provider without a version constraint:

```hcl
provider "<REPLACE_ME>" {}
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