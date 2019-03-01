## Intro

This is an experiemental terraform provider for managing your Mikrotik's router's DNS. It allows you to create records via the routeos API.

## Examples

Creating a new DNS record

```terraform
provider "mikrotik" {
  host = "http://router:8728" # Or set MIKROTIK_HOST
  username = "username of api user" # Or set MIKROTIK_USER
  password = "xxxxxx" #  # Or set MIKROTIK_PASSWORD
}

resource "mikrotik_dns_record" "www" {
    name = "router"
    address = "192.168.88.1"
    ttl = 300
}
```

## Todo
- [ ] Add Travis test suite
- [ ] Add more rigorious Terraform tests
- [ ] Resource reading needs to be more robust so the terraform plan does not think it needs to recreate everything when really the state is fine.
