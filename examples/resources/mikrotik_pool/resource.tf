resource "mikrotik_pool" "pool" {
  name    = "pool-name"
  ranges  = "172.16.0.6-172.16.0.12"
  comment = "ip pool with range specified"
}
