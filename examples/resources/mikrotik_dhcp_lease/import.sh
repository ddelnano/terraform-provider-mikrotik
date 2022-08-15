# The resource ID (*19) is a MikroTik's internal id.
# It can be obtained via CLI:
# [admin@MikroTik] /ip dhcp-server lease> :put [find where address=10.0.1.254]
# *19
terraform import mikrotik_dhcp_lease.file_server '*19'
