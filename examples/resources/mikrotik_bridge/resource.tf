resource "mikrotik_bridge" "bridge" {
  name           = "default_bridge"
  fast_forward   = true
  vlan_filtering = false
  comment        = "Default bridge"
}
