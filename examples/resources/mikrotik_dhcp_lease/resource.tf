resource "mikrotik_dhcp_lease" "file_server" {
  address    = "192.168.88.1"
  macaddress = "11:22:33:44:55:66"
  comment    = "file server"
  blocked    = "false"
}
