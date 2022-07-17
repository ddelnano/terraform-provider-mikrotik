# mikrotik_dhcp_server
Creates a DHCP server on the mikrotik device

## Example Usage

```hcl
resource "mikrotik_dhcp_server" "record" {
  address_pool  = "lan-pool"
  authoritative = "yes"
  disabled      = false
  interface     = "ether2"
  name          = "lan-dhcp-server"
}
```

## Argument Reference
* name - (Required) Reference name.
* add_arp - (Optional) Whether to add dynamic ARP entry.
* address_pool - (Optional) IP pool, from which to take IP addresses for the clients.
* authoritative - (Optional) The way how server responds to DHCP requests.
* disabled - (Optional) Whether this instance of DHCP server should be disabled.
* interface - (Optional) Interface on which server will be running.
* lease_script - (Optional) Script that will be executed after lease is assigned or de-assigned.

## Attributes Reference

## Import Reference

```bash
terraform import mikrotik_dhcp_server.lan "my-dhcp-server"
```
