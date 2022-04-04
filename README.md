# Terraform Provider for Ionos Developer

The IonosDeveloper provider gives the ability to configure DNS records using the IONOS Developer API.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.13.x+
- [Go](https://golang.org/doc/install) 1.17 (to build the provider plugin)

  **NOTE:** In order to use a specific version of this provider, please include the following block at the beginning of your terraform config files [details](https://www.terraform.io/docs/configuration/terraform.html#specifying-a-required-terraform-version):

```terraform
provider "ionosdeveloper" {
  version = "~> 0.0.1"
}
```

## Building The Provider

Clone repository to: `$GOPATH/src/github.com/ionos-developer/terraform-provider-ionosdeveloper`

```sh
$ mkdir -p $GOPATH/src/github.com/ionos-developer; cd $GOPATH/src/github.com/ionos-developer
$ git clone git@github.com:ionos-developer/terraform-provider-ionosdeveloper
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/ionos-developer/terraform-provider-ionosdeveloper
$ make build
```

## Using the provider

See the [IonosDeveloper Provider documentation](https://registry.terraform.io/providers/ionos-developer/ionosdeveloper/latest/docs) to get started using the IonosDeveloper provider.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.17+ is _required_). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-ionosdeveloper
...
```

## Testing the Provider

### What Are We Testing?

The purpose of our acceptance tests is to **provision** resources containing all the available arguments, followed by **updates** on all arguments that allow this action. Beside the provisioning part, **data-sources** with all possible arguments are also tested.

All tests are integrated into [github actions](https://github.com/ionos-developer/terraform-provider-ionosdeveloper/actions) that run daily and are also run manually before any release.

### How to Run Tests Locally

In order to test the provider, you can simply run:

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run:

```sh
$ make testacc TAGS=all
```

#### Test Tags

Tests can also be run for a batch of resources using tags.

_Example of running dns tests:_

```sh
$ make testacc TAGS=dns
```

<details> <summary title="Click to toggle">See more details about <b>test tags</b></summary>

**Build tags** are named as follows:

- `dns` - all **dns** tests

</details>
