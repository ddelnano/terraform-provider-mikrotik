resource "mikrotik_bpg_instance" "instance" {
  name      = "bgp-instance-name"
  as        = 65533
  router_id = "172.21.16.20"
  comment   = "test comment"
}

resource "mikrotik_bgp_peer" "peer" {
  name           = "bgp-peer-name"
  remote_as      = 65533
  remote_address = "172.21.16.20"
  instance       = mikrotik_bgp_instance.instance.name
}
