resource "mikrotik_bgp_instance" "instance" {
  name      = "bgp-instance-name"
  as        = 65533
  router_id = "172.21.16.20"
  comment   = "test comment"
}
