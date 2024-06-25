# The ID argument (*2) is a MikroTik's internal id.
# It can be obtained via CLI:
#
# [admin@MikroTik] /ip dns static> :put [find where address="192.168.88.1/24"]
# *2
terraform import mikrotik_dns_record.record "*2"
