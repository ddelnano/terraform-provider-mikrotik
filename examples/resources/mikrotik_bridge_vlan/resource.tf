resource "mikrotik_bridge" "default" {
  name = "main"
}

resource "mikrotik_bridge_vlan" "testacc" {
  bridge   = mikrotik_bridge.default.name
  tagged   = ["ether2", "vlan30"]
  untagged = ["ether3"]
  vlan_ids = [10, 30]
}
