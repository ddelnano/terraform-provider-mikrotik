# The ID argument (*19) is a MikroTik's internal id.
# It can be obtained via CLI:
#
# [admin@MikroTik] /ipv6 address> :put [find where address="192.168.88.1/24"]
# *19
terraform import mikrotik_ipv6_address.lan *19
