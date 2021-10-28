## mikrotik_ip_address

Assignes an IP address to an interface

### Example Usage

```hcl
resource "mikrotik_ip_address" "lan" {
  address = "192.168.88.1/24"
  comment = "LAN Network"
  interface = "ether1"
}
```

### Argument Reference
* address - (Required) The IP address and netmask of the interface using slash notation
* comment - (Optional) The comment for the IP address assignment
* disabled - (Optional) Whether to disable IP address (true\false)
* interface - (Required) The interface on which the IP address is assigned

### Attributes Reference

### Import Reference

```bash
terraform import mikrotik_ip_address.lan *19
```

Last argument (*19) is a mikrotik internal id which can be obtained via CLI:

```
[admin@MikroTik] /ip address> :put [find where address="192.168.88.1/24"]
*19
```
