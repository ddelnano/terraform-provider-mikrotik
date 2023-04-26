# Mikrotik provider for Terraform 

## Intro

This is a terraform provider for managing resources on your RouterOS device. To see what resources and data sources are supported, please see the [documentation](https://registry.terraform.io/providers/ddelnano/mikrotik/latest/docs) on the terraform registry.

## Support

You can discuss any issues you have or feature requests in [Discord](https://discord.gg/ZpNq8ez).

## Donations

If you get value out this project and want to show your support you can find me on [patreon](https://www.patreon.com/ddelnano).

## Building provider locally

Requirements:
* [Go](https://go.dev/doc/install) >= 1.18
* [Terraform]() >= 0.14

To build the provider with `make`:
```shell
$ make build
```
which creates a `terraform-provider-mikrotik` binary in repository's root folder.

or build with `go` compiler:
```shell
$ go build -o terraform-provider-mikrotik
```

To use locally built provider, Terraform should be aware of its binary.

It could be done with custom CLI config file:
```hcl
# custom.tfrc

provider_installation {
    dev_overrides {
        "ddelnano/mikrotik" = "/path/to/clones/repository/terraform-provider-mikrotik"
    }

    direct {}
}
```
The [dev_overrides](https://developer.hashicorp.com/terraform/cli/config/config-file#development-overrides-for-provider-developers) section is available since Terraform `0.14`.

Finally, tell Terraform CLI to use custom confiuration by exporting environment variable:
```shell
$ export TF_CLI_CONFIG_FILE=path/to/custom.tfrc
```

**NOTE**: with `dev_overrides` it is not possible to run `terraform init` (see [official docs](https://developer.hashicorp.com/terraform/cli/config/config-file#development-overrides-for-provider-developers)) so you should immediately use `terraform plan` and `terraform apply` without initializing.

## Contributing

### Dependencies
- RouterOS. See which versions are supported by what is tested in [CI](.github/workflows/continuous-integration.yml)
- Terraform 0.12+

### Testing

The provider is tested with Terraform's acceptance testing framework. As long as you have a RouterOS device you should be able to run them. Please be aware it will create resources on your device! Code that is accepted by the project will not be destructive for anything existing on your router but be careful when changing test code!

In order to run the tests you will need to set the following environment variables:
```bash
export MIKROTIK_HOST=router-hostname:8728
export MIKROTIK_USER=username
# Please be aware this will put your password in your bash history and is not safe
export MIKROTIK_PASSWORD=password
```

After those environment variables are set you can run the tests with the following command:
```bash
make testacc
```

> :warning:
>
> For RouterOS versions `6.x` you have to set environment variable `LEGACY_BGP_SUPPORT="true"` before running tests.

### Testing without MikroTik hardware

If you do not have MikroTik hardware or virtual machine with pre-installed RouterOS, you still have a way to run tests locally.

To make this happen, install [Docker](https://www.docker.com) on your developer machine, and run from the root of the repository:
```sh
$ make routeros
```
It will start RouterOS container locally and make its API server available at `127.0.0.1:8728`

Just export connection settings
```sh
export MIKROTIK_HOST=127.0.0.1:8728
export MIKROTIK_USER=admin
export MIKROTIK_PASSWORD=""
```

and you are ready to run tests with
```sh
$ make test
```

You can use specific RouterOS version by passing `ROUTEROS_VERSION` argument
```sh
$ make routeros ROUTEROS_VERSION="6.49beta54"
```

or even
```sh
$ make routeros ROUTEROS_VERSION=latest
```

To cleanup everything, just run:
```sh
$ make routeros clean
```
