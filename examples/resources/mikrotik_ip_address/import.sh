# The ID argument (*19) is a MikroTik's internal id.
# It can be obtained via CLI:
#
# [admin@MikroTik] /ip address> :put [find where address="192.168.88.1/24"]
# *19
terraform import mikrotik_ip_address.lan '*19'
