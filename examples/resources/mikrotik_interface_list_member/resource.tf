resource "mikrotik_interface_list" "lan" {
  name = "lan"
}

resource "mikrotik_interface_list_member" "lan" {
  interface = "ether2"
  list      = mikrotik_interface_list.lan.name
}
