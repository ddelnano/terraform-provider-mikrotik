package client

import (
	"github.com/ddelnano/terraform-provider-mikrotik/client/types"
)

type Scheduler struct {
	Id        string                 `mikrotik:".id"`
	Name      string                 `mikrotik:"name"`
	OnEvent   string                 `mikrotik:"on-event"`
	StartDate string                 `mikrotik:"start-date"`
	StartTime string                 `mikrotik:"start-time"`
	Interval  types.MikrotikDuration `mikrotik:"interval"`
}

func (s *Scheduler) ActionToCommand(action Action) string {
	return map[Action]string{
		Add:    "/system/scheduler/add",
		Find:   "/system/scheduler/print",
		List:   "/system/scheduler/print",
		Update: "/system/scheduler/set",
		Delete: "/system/scheduler/remove",
	}[action]
}

func (s *Scheduler) IDField() string {
	return ".id"
}

func (s *Scheduler) ID() string {
	return s.Id
}

func (s *Scheduler) SetID(id string) {
	s.Id = id
}

func (s *Scheduler) FindField() string {
	return "name"
}

func (s *Scheduler) FindFieldValue() string {
	return s.Name
}

func (s *Scheduler) DeleteField() string {
	return "numbers"
}

func (s *Scheduler) DeleteFieldValue() string {
	return s.Id
}

func (client Mikrotik) FindScheduler(name string) (*Scheduler, error) {
	res, err := client.Find(&Scheduler{Name: name})
	if err != nil {
		return nil, err
	}

	return res.(*Scheduler), nil
}

func (client Mikrotik) DeleteScheduler(name string) error {
	return client.Delete(&Scheduler{Id: name})
}

// AddScheduler is an alias to CreateScheduler
func (client Mikrotik) AddScheduler(s *Scheduler) (*Scheduler, error) {
	return client.CreateScheduler(s)
}

func (client Mikrotik) CreateScheduler(s *Scheduler) (*Scheduler, error) {
	res, err := client.Add(s)
	if err != nil {
		return nil, err
	}

	return res.(*Scheduler), nil
}

func (client Mikrotik) UpdateScheduler(s *Scheduler) (*Scheduler, error) {
	res, err := client.Update(s)
	if err != nil {
		return nil, err
	}

	return res.(*Scheduler), nil
}
