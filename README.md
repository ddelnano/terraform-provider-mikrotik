# Mikrotik provider for Terraform 

## Intro

This is a terraform provider for managing resources on your RouterOS device. It currently has support for the following RouterOS concepts:
- DNS records
- Scripts

## Examples

See the [pages directory](modules/reference/pages) for the documentation on each resource.

## Support

You can discuss any issues you have or feature requests in [Discord](https://discord.gg/ZpNq8ez).

## Donations

If you get value out this project and want to show your support you can find me on [patreon](https://www.patreon.com/ddelnano).

## Contributing

### Dependencies
- go 1.11+ to build the provider
- RouterOS v6.45.2+ (It may work on older versions but it is untested)

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
