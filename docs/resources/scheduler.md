# mikrotik_scheduler

Creates a Mikrotik scheduler.

## Example Usage

```hcl
resource "mikrotik_scheduler" "scheduler" {
  name = "scheduler-name"
  on_event = "scheduler-to-execute"
  # Run every 5 mins
  interval = 300
```

## Argument Reference
* name - (Required) The name of scheduler.
* on_event - (Required) The name of the script to run
* interval - (Optiona) The interval between two script executions, if time interval is set to zero, the script is only executed at its start time, otherwise it is executed repeatedly at the time interval is specified

See the https://wiki.mikrotik.com/wiki/Manual:System/Scheduler[mikrotik docs] for more details

## Attributes Reference

## Import Reference

```bash
terraform import mikrotik_scheduler.scheduler scheduler-name
```
