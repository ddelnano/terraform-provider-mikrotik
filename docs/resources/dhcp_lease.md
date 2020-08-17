## mikrotik_dhcp_lease

Creates a DHCP lease on the mikrotik device

### Example Usage

```hcl
resource "mikrotik_dhcp_lease" "file_server" {
  address = "192.168.88.1"
  macaddress = "11:22:33:44:55:66"
  comment = "file server"
}
```

### Argument Reference
* address - (Required) The IP address of the DHCP lease to be created
* macaddress - (Required) The MAC addreess of the DHCP lease to be created
* comment - (Optional) The comment of the DHCP lease to be created

### Attributes Reference

### Import Reference

```bash
terraform import mikrotik_dhcp_lease.file_server *19
```

Last argument (*19) is a mikrotik internal id which can be obtained via CLI:

```
[admin@MikroTik] /ip dhcp-server lease> :put [find where address=10.0.1.254]
*19
```
