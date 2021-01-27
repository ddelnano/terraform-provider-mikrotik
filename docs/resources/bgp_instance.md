# mikrotik_pool

Creates a Mikrotik [BGP Instance](https://wiki.mikrotik.com/wiki/Manual:Routing/BGP#Instance).

## Example Usage

```hcl
resource "mikrotik_bgp_instance" "instance" {
  name      = "bgp-instance-name"
  as        = 65533
  router_id = 172.21.16.20
  comment   = "test comment"
}
```

## Argument Reference
* name - (Required) The name of the BGP instance.
* as - (Required) The 32-bit BGP autonomous system number. Must be a value within 0 to 4294967295.
* router_id - (Required) BGP Router ID (for this instance). If set to 0.0.0.0, BGP will use one of router's IP addresses.
* routing_table - (Optional, Default: `""`)	Name of routing table this BGP instance operates on.
* client_to_client_reflection - (Optional, Default: `true`) The comment of the IP Pool to be created
* comment - (Optional) The comment of the BGP instance to be created.
* confederation - (Optional) Autonomous system number that identifies the [local] confederation as a whole. Must be a value within 0 to 4294967295.
* confederation_peers - (Optional) List of AS numbers internal to the [local] confederation. For example:  `"10,20,30-50"`
* disabled - (Optional, Default: `true`) Whether instance is disabled.
* ignore_as_path_len - (Optional, Default: `false`) Whether to ignore AS_PATH attribute in BGP route selection algorithm.
* out_filter - (Optional, Default: `""`) Output routing filter chain used by all BGP peers belonging to this instance.
* redistribute_connected - (Optional, Default: `false`) If enabled, this BGP instance will redistribute the information about connected routes.
* redistribute_ospf - (Optional, Default: `false`) If enabled, this BGP instance will redistribute the information about routes learned by OSPF.
* redistribute_other-bgp - (Optional, Default: `false`) If enabled, this BGP instance will redistribute the information about routes learned by other BGP instances.
* redistribute_rip - (Optional, Default: `false`)	If enabled, this BGP instance will redistribute the information about routes learned by RIP.
* redistribute_static - (Optional, Default: `false`)	If enabled, the router will redistribute the information about static routes added to its routing database.


## Attributes Reference

## Import Reference

```bash
# import with name of bgp instance
terraform import mikrotik_bgp_instance.instance bgp-instance-name
```
