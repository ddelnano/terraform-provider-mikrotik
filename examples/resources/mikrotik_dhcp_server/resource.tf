resource "mikrotik_pool" "bar" {
  name    = "dhcp-pool"
  ranges  = "10.10.10.100-10.10.10.200"
  comment = "Home devices"
}

resource "mikrotik_dhcp_server" "default" {
  address_pool  = mikrotik_pool.bar.name
  authoritative = "yes"
  disabled      = false
  interface     = "ether2"
  name          = "main-dhcp-server"
}
