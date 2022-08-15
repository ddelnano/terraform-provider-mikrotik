resource "mikrotik_script" "script" {
  name  = "script-name"
  owner = "admin"
  policy = [
    "ftp",
    "reboot",
  ]
  source = <<EOF
:put testing
EOF
}
