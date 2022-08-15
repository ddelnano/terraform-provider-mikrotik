resource "mikrotik_ipv6_address" "lan" {
  address   = "2001::1/64"
  comment   = "LAN Network"
  interface = "ether1"
}
