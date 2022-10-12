# mikrotik_bridge (Resource)
Manages a bridge resource on remote MikroTik device.

## Example Usage
```terraform
resource "mikrotik_bridge" "bridge" {
  name           = "default_bridge"
  fast_forward   = true
  vlan_filtering = false
  comment        = "Default bridge"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name of the bridge interface

### Optional

- `comment` (String) Short description of the interface.
- `fast_forward` (Boolean) Special and faster case of FastPath which works only on bridges with 2 interfaces (enabled by default only for new bridges). Default: `true`.
- `vlan_filtering` (Boolean) Globally enables or disables VLAN functionality for bridge.

### Read-Only

- `id` (String) The ID of this resource.

## Import
Import is supported using the following syntax:
```shell
# import with name of bridge
terraform import mikrotik_bridge.bridge <bridge_name>
```