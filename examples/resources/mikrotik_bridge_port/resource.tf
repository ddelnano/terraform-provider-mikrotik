resource "mikrotik_bridge" "bridge" {
  name           = "default_bridge"
  fast_forward   = true
  vlan_filtering = false
  comment        = "Default bridge"
}

resource mikrotik_bridge_port "eth2port" {
  bridge    = mikrotik_bridge.bridge.name
  interface = "ether2"
  pvid      = 10
  comment   = "bridge port"
}
