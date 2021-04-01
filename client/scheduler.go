package client

import (
	"fmt"
	"log"
)

type Scheduler struct {
	Id        string `mikrotik:".id"`
	Name      string
	OnEvent   string `mikrotik:"on-event"`
	StartDate string `mikrotik:"start-date"`
	StartTime string `mikrotik:"start-time"`
	Interval  int    `mikrotik:"interval,ttlToSeconds"`
}

func (client Mikrotik) FindScheduler(name string) (*Scheduler, error) {
	c, err := client.getMikrotikClient()
	cmd := []string{"/system/scheduler/print", "?name=" + name}
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)

	log.Printf("[DEBUG] Found scheduler from mikrotik api %v", r)
	scheduler := &Scheduler{}
	err = Unmarshal(*r, scheduler)

	if err != nil {
		return nil, err
	}

	if scheduler.Name == "" {
		return nil, NewNotFound(fmt.Sprintf("scheduler `%s` not found", name))
	}
	return scheduler, err
}

func (client Mikrotik) DeleteScheduler(name string) error {
	c, err := client.getMikrotikClient()

	scheduler, err := client.FindScheduler(name)

	if err != nil {
		return err
	}
	cmd := []string{"/system/scheduler/remove", "=numbers=" + scheduler.Id}
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)
	log.Printf("[DEBUG] Remove scheduler from mikrotik api %v", r)

	return err
}

func (client Mikrotik) CreateScheduler(name string, onEvent string, interval int) (*Scheduler, error) {
	c, err := client.getMikrotikClient()

	nameArg := fmt.Sprintf("=name=%s", name)
	onEventArg := fmt.Sprintf("=on-event=%s", onEvent)
	intervalArg := fmt.Sprintf("=interval=%d", interval)
	cmd := []string{
		"/system/scheduler/add",
		nameArg,
		onEventArg,
		intervalArg,
	}
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)
	log.Printf("[DEBUG] /system/scheduler/add returned %v", r)

	if err != nil {
		return nil, err
	}

	return client.FindScheduler(name)
}

func (client Mikrotik) UpdateScheduler(name, onEvent string, interval int) (*Scheduler, error) {
	c, err := client.getMikrotikClient()

	if err != nil {
		return nil, err
	}

	scheduler, err := client.FindScheduler(name)

	if err != nil {
		return scheduler, err
	}

	nameArg := fmt.Sprintf("=numbers=%s", scheduler.Id)
	intervalArg := fmt.Sprintf("=interval=%d", interval)
	onEventArg := fmt.Sprintf("=on-event=%s", onEvent)
	cmd := []string{
		"/system/scheduler/set",
		nameArg,
		intervalArg,
		onEventArg,
	}
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	_, err = c.RunArgs(cmd)

	if err != nil {
		return scheduler, err
	}

	return client.FindScheduler(name)
}
