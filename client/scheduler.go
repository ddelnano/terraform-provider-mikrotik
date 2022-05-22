package client

import (
	"fmt"
	"log"
)

//go:generate gen
type Scheduler struct {
	Id        string `mikrotik:".id" gen:"-,mikrotikID"`
	Name      string `mikrotik:"name" gen:"name,id,required"`
	OnEvent   string `mikrotik:"on-event" gen:"on_event,required"`
	StartDate string `mikrotik:"start-date" gen:"start_date,computed"`
	StartTime string `mikrotik:"start-time" gen:"start_time,computed"`
	Interval  int    `mikrotik:"interval,ttlToSeconds" gen:"interval,optional"`
}

func (client Mikrotik) FindScheduler(name string) (*Scheduler, error) {
	c, err := client.getMikrotikClient()

	if err != nil {
		return nil, err
	}

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

	if err != nil {
		return err
	}

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

func (client Mikrotik) AddScheduler(s *Scheduler) (*Scheduler, error) {
	return client.CreateScheduler(s)
}

func (client Mikrotik) CreateScheduler(s *Scheduler) (*Scheduler, error) {
	c, err := client.getMikrotikClient()

	if err != nil {
		return nil, err
	}

	cmd := Marshal("/system/scheduler/add", s)

	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)
	log.Printf("[DEBUG] /system/scheduler/add returned %v", r)

	if err != nil {
		return nil, err
	}

	return client.FindScheduler(s.Name)
}

func (client Mikrotik) UpdateScheduler(s *Scheduler) (*Scheduler, error) {
	c, err := client.getMikrotikClient()

	if err != nil {
		return nil, err
	}

	scheduler, err := client.FindScheduler(s.Name)

	if err != nil {
		return scheduler, err
	}

	cmd := Marshal("/system/scheduler/set", s)

	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	_, err = c.RunArgs(cmd)

	if err != nil {
		return scheduler, err
	}

	return client.FindScheduler(s.Name)
}
