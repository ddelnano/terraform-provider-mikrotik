# The ID argument (*17) is a MikroTik's internal id.
# It can be obtained via CLI:
#
# [admin@MikroTik] /ip pool> :put [ find where name=pool-name]
# *17
terraform import mikrotik_pool.pool '*17'
