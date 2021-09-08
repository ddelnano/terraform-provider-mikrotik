# mikrotik_bgp_peer

Creates a Mikrotik [BGP Peer](https://wiki.mikrotik.com/wiki/Manual:Routing/BGP#Peer).

This resource will not be supported in RouterOS v7+. Mikrotik has deprecated the underlying commands so future BGP support will need new resources created (See #52 for status of this work).

## Example Usage

```hcl
resource "mikrotik_bpg_instance" "instance" {
  [...]
}

resource "mikrotik_bgp_peer" "peer" {
  name = "bgp-peer-name"
  remote_as = 65533
  remote_address = "172.21.16.20"
  instance = mikrotik_bgp_instance.instance.name
}
```

## Argument Reference
* name - (Required) The name of the BGP peer.
* remote_as - (Required) The 32-bit AS number of the remote peer.
* remote_address - (Required) The address of the remote peer
* instance - (Required)	The name of the instance this peer belongs to.
See [Mikrotik bgp instance resource](https://github.com/ddelnano/terraform-provider-mikrotik/blob/master/docs/resources/bgp_instance.md).
* address_families - (Optional, Default: `"ip"`) The list of address families about which this peer will exchange routing information.
* ttl - (Optional, Default: `"default"`) Time To Live, the hop limit for TCP connection. This is a `string` field that can be 'default' or '0'-'255'.
* default_originate - (Optional, Default: `"never"`) The comment of the BGP peer to be created.
* hold_time - (Optional, Default: "3m") Specifies the BGP Hold Time value to use when negotiating with peer
* nexthop_choice - (Optional, Default: "default") Affects the outgoing NEXT_HOP attribute selection, either:  'default', 'force-self', or 'propagate'
* disabled - (Optional, Default: `true`) Whether peer is disabled.
* comment - (Optional) The comment of the BGP peer to be created.
* out_filter - (Optional) The name of the routing filter chain that is applied to the outgoing routing information. 
* in_filter - (Optional) The name of the routing filter chain that is applied to the incoming routing information.
* allow_as_in - (Optional) How many times to allow own AS number in AS-PATH, before discarding a prefix.
* as_override - (Optional, Default: `false`) If set, then all instances of remote peer's AS number in BGP AS PATH attribute are replaced with local AS number before sending route update to that peer.
* cisco_vpls_nlri_len_fmt - (Optional) VPLS NLRI length format type.
* max_prefix_limit - (Optional) Maximum number of prefixes to accept from a specific peer.
* keepalive_time - (Optional)
* max_prefix_restart_time - (Optional) Minimum time interval after which peers can reestablish BGP session.
* multihop - (Optional, Default: `false`)	Specifies whether the remote peer is more than one hop away.
* passive - (Optional, Default: `false`) Name of the routing filter chain that is applied to the outgoing routing information.
* remote_port - (Optional) Remote peers port to establish tcp session.
* remove_private_as - (Optional, Default: `false`) If set, then BGP AS-PATH attribute is removed before sending out route update if attribute contains only private AS numbers.
* route_refelct - (Optional, Default: `false`) Specifies whether this peer is route reflection client.
* tcp_md5_key - (Optional) Key used to authenticate the connection with TCP MD5 signature as described in RFC 2385.
* update_source - (Optional) If address is specified, this address is used as the source address of the outgoing TCP connection.
* use_bfd - (Optional, Default: `false`) Whether to use BFD protocol for fast state detection.


## Attributes Reference

## Import Reference

```bash
# import with name of bgp peer 
terraform import mikrotik_bgp_peer.peer bgp-peer-name
```
