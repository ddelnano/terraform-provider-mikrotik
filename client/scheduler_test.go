package client

import (
	"testing"

	"github.com/ddelnano/terraform-provider-mikrotik/client/types"
	"github.com/stretchr/testify/require"
)

func TestCreateUpdateDeleteAndFindScheduler(t *testing.T) {
	c := NewClient(GetConfigFromEnv())

	schedulerName := "scheduler_" + RandomString()
	onEvent := "onevent"
	interval := 0
	expectedScheduler := &Scheduler{
		Name:     schedulerName,
		OnEvent:  onEvent,
		Interval: types.MikrotikDuration(interval),
	}
	scheduler, err := c.CreateScheduler(expectedScheduler)
	require.NoError(t, err)
	require.NotNil(t, scheduler)

	expectedScheduler.Id = scheduler.Id
	expectedScheduler.StartDate = scheduler.StartDate
	expectedScheduler.StartTime = scheduler.StartTime

	require.Equal(t, expectedScheduler, scheduler)

	// update and reassert
	expectedScheduler.OnEvent = "test"
	scheduler, err = c.UpdateScheduler(expectedScheduler)
	require.Equal(t, expectedScheduler, scheduler)

	err = c.DeleteScheduler(schedulerName)
	require.NoError(t, err)
}

func TestFindScheduler_onNonExistantScript(t *testing.T) {
	c := NewClient(GetConfigFromEnv())

	name := "scheduler does not exist"
	_, err := c.FindScheduler(name)

	if _, ok := err.(*NotFound); !ok {
		t.Errorf("Expecting to receive NotFound error for scheduler `%s`, instead error was nil.", name)
	}
}
