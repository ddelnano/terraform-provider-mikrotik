# Mikrotik Provider

The mikrotik provider is used to interact with the resources supported by RouterOS.
The provider needs to be configured with the proper credentials before it can be used.

## Requirements

* RouterOS v6.45.2+ (It may work with other versions but it is untested against other versions!)

## Using the provider

This provider is on the terraform registry so you only need to reference it in your terraform code (example below).

## Example Usage

```hcl
# Configure the mikrotik Provider
provider "mikrotik" {
  host = "hostname-of-server:8728"     # Or set MIKROTIK_HOST environment variable
  username = "<username>"              # Or set MIKROTIK_USER environment variable
  password = "<password>"              # Or set MIKROTIK_PASSWORD environment variable
  tls = true|false                     # Or set MIKROTIK_TLS environment variable
  ca_certificate = "/path/to/ca/certificate.pem"    # Or set MIKROTIK_CA_CERTIFICATE environment variable
  insecure = true|false                # Or set MIKROTIK_INSECURE environment variable
}
```
