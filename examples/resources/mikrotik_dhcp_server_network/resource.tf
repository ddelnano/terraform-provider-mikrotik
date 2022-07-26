resource "mikrotik_dhcp_server_network" "default" {
  address    = "192.168.100.0/24"
  netmask    = "0" # use mask from address
  gateway    = "192.168.100.1"
  dns_server = "192.168.100.2"
  comment    = "Default DHCP server network"
}
