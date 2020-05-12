package client

import (
	"reflect"
	"testing"
)

func TestCreateDeleteAndFindScheduler(t *testing.T) {
	c := NewClient(GetConfigFromEnv())

	schedulerName := "scheduler"
	onEvent := "onevent"
	interval := 0
	expectedScheduler := Scheduler{
		Name:    schedulerName,
		OnEvent: onEvent,
	}
	scheduler, err := NewClient(GetConfigFromEnv()).CreateScheduler(
		schedulerName,
		onEvent,
		interval,
	)

	if err != nil || scheduler == nil {
		t.Errorf("Error creating a scheduler with: %v and value: %v", err, scheduler)
	}

	expectedScheduler.Id = scheduler.Id
	expectedScheduler.StartDate = scheduler.StartDate
	expectedScheduler.StartTime = scheduler.StartTime

	if !reflect.DeepEqual(*scheduler, expectedScheduler) {
		t.Errorf("The scheduler does not match what we expected. actual: %v expected: %v", *scheduler, expectedScheduler)
	}

	err = c.DeleteScheduler(schedulerName)

	if err != nil {
		t.Errorf("Error deleting a scheduler with: %v", err)
	}

	scheduler, err = c.FindScheduler(schedulerName)

	if err != nil || scheduler != nil {
		t.Errorf("Scheduler (%v) was not deleted: %v", scheduler, err)
	}
}
