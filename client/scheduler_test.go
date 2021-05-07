package client

import (
	"reflect"
	"testing"
)

func TestCreateUpdateDeleteAndFindScheduler(t *testing.T) {
	c := NewClient(GetConfigFromEnv())

	schedulerName := "scheduler"
	onEvent := "onevent"
	interval := 0
	expectedScheduler := &Scheduler{
		Name:     schedulerName,
		OnEvent:  onEvent,
		Interval: interval,
	}
	scheduler, err := c.CreateScheduler(expectedScheduler)

	if err != nil || scheduler == nil {
		t.Errorf("Error creating a scheduler with: %v and value: %v", err, scheduler)
	}

	expectedScheduler.Id = scheduler.Id
	expectedScheduler.StartDate = scheduler.StartDate
	expectedScheduler.StartTime = scheduler.StartTime

	if !reflect.DeepEqual(scheduler, expectedScheduler) {
		t.Errorf("The scheduler does not match what we expected. actual: %v expected: %v", scheduler, expectedScheduler)
	}

	// update and reassert
	expectedScheduler.OnEvent = "test"
	scheduler, err = c.UpdateScheduler(expectedScheduler)

	if !reflect.DeepEqual(scheduler, expectedScheduler) {
		t.Errorf("The updated scheduler does not match what we expected. actual: %v expected: %v", scheduler, expectedScheduler)
	}

	err = c.DeleteScheduler(schedulerName)

	if err != nil {
		t.Errorf("Error deleting a scheduler with: %v", err)
	}
}

func TestFindScheduler_onNonExistantScript(t *testing.T) {
	c := NewClient(GetConfigFromEnv())

	name := "scheduler does not exist"
	_, err := c.FindScheduler(name)

	if _, ok := err.(*NotFound); !ok {
		t.Errorf("Expecting to receive NotFound error for scheduler `%s`, instead error was nil.", name)
	}
}
