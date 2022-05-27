## mikrotik_ipv6_address

Assigns an IP address to an interface

### Example Usage

```hcl
resource "mikrotik_ipv6_address" "lan" {
  address = "2001::1/64"
  comment = "LAN Network"
  interface = "ether1"
}
```

### Argument Reference
* address - (Required) The IPv6 address and prefix length of the interface using slash notation.
* advertise - (Optional) Whether to enable stateless address configuration. The prefix of that address is automatically advertised to hosts using ICMPv6 protocol. The option is set by default for addresses with prefix length 64.
* comment - (Optional) The comment for the IPv6 address assignment.
* disabled - (Optional) Whether to disable IPv6 address (true\false).
* eui_64 - (Optional) Whether to calculate EUI-64 address and use it as last 64 bits of the IPv6 address.
* from_pool - (Optional) Name of the pool from which prefix will be taken to construct IPv6 address taking last part of the address from address property.
* interface - (Required) The interface on which the IPv6 address is assigned.
* no_dad - (Optional) If set indicates that address is anycast address and Duplicate Address Detection should not be performed.

### Attributes Reference

### Import Reference

```bash
terraform import mikrotik_ipv6_address.lan *19
```

Last argument (*19) is a mikrotik internal id which can be obtained via CLI:

```
[admin@MikroTik] /ipv6 address> :put [find where address="2001::1/64"]
*19
```
