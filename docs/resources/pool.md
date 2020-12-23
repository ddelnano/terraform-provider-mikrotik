# mikrotik_pool

Creates a Mikrotik [IP Pool](https://wiki.mikrotik.com/wiki/Manual:IP/Pools).

## Example Usage

```hcl
resource "mikrotik_pool" "pool" {
  name = "pool-name"
  ranges = "172.16.0.6-172.16.0.12"
  comment = "ip pool with range specified"
}
```

## Argument Reference
* name - (Required) The name of IP pool.
* ranges - (Required) The IP range(s) of the pool. Multiple ranges can be specified, separated by commas (e.g. 172.16.0.6-172.16.0.12,172.16.0.50-172.16.0.60)
* comment - (Optional) The comment of the IP Pool to be created

## Attributes Reference

## Import Reference

```bash
terraform import mikrotik_pool.pool *17
```

Last argument (*17) is a mikrotik internal id which can be obtained via CLI:

```
[admin@MikroTik] /ip pool> :put [ find where name=pool-name]
*17
```
