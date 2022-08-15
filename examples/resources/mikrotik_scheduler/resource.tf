resource "mikrotik_scheduler" "scheduler" {
  name     = "scheduler-name"
  on_event = "scheduler-to-execute"
  # Run every 5 mins
  interval = 300
}
