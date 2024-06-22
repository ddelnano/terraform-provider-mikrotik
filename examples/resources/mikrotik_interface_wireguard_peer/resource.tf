resource "mikrotik_interface_wireguard" "default" {
  name    = "wireguard-interface"
  comment = "new interface"
}

resource "mikrotik_interface_wireguard_peer" "default" {
  interface       = mikrotik_interface_wireguard.default.name
  public_key      = "v/oIzPyFm1FPHrqhytZgsKjU7mUToQHLrW+Tb5e601M="
  comment         = "peer-1"
  allowed_address = "0.0.0.0/0"
}
