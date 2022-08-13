resource "mikrotik_vlan_interface" "default" {
  interface = "ether2"
  mtu       = 1500
  name      = "vlan-20"
  vlan_id   = 20
}
