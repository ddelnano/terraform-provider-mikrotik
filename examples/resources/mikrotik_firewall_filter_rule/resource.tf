resource "mikrotik_firewall_filter_rule" "https" {
  action             = "accept"
  chain              = "forward"
  comment            = "Web access to local HTTP server"
  connection_state   = ["new"]
  dst_port           = "443"
  in_interface       = "ether1"
  in_interface_list  = "local_lan"
  out_interface_list = "ether3"
  protocol           = "tcp"
}
