resource "mikrotik_ip_address" "lan" {
  address   = "192.168.88.1/24"
  comment   = "LAN Network"
  interface = "ether1"
}
