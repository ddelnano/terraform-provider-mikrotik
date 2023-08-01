package client

import (
	"github.com/ddelnano/terraform-provider-mikrotik/client/types"
	"github.com/go-routeros/routeros"
)

type Scheduler struct {
	Id        string                 `mikrotik:".id"`
	Name      string                 `mikrotik:"name"`
	OnEvent   string                 `mikrotik:"on-event"`
	StartDate string                 `mikrotik:"start-date"`
	StartTime string                 `mikrotik:"start-time"`
	Interval  types.MikrotikDuration `mikrotik:"interval"`
}

var _ Resource = (*Scheduler)(nil)

func (b *Scheduler) ActionToCommand(a Action) string {
	return map[Action]string{
		Add:    "/system/scheduler/add",
		Find:   "/system/scheduler/print",
		Update: "/system/scheduler/set",
		Delete: "/system/scheduler/remove",
	}[a]
}

func (b *Scheduler) IDField() string {
	return ".id"
}

func (b *Scheduler) ID() string {
	return b.Id
}

func (b *Scheduler) SetID(id string) {
	b.Id = id
}

func (b *Scheduler) AfterAddHook(r *routeros.Reply) {
	b.Id = r.Done.Map["ret"]
}

func (b *Scheduler) FindField() string {
	return "name"
}

func (b *Scheduler) FindFieldValue() string {
	return b.Name
}

func (b *Scheduler) DeleteField() string {
	return "numbers"
}

func (b *Scheduler) DeleteFieldValue() string {
	return b.Id
}

// typed wrappers
func (client Mikrotik) AddScheduler(s *Scheduler) (*Scheduler, error) {
	return client.CreateScheduler(s)
}

func (client Mikrotik) CreateScheduler(s *Scheduler) (*Scheduler, error) {
	r, err := client.Add(s)
	if err != nil {
		return nil, err
	}

	return r.(*Scheduler), nil
}

func (client Mikrotik) UpdateScheduler(s *Scheduler) (*Scheduler, error) {
	r, err := client.Update(s)
	if err != nil {
		return nil, err
	}

	return r.(*Scheduler), nil
}

func (client Mikrotik) FindScheduler(name string) (*Scheduler, error) {
	r, err := client.Find(&Scheduler{Name: name})
	if err != nil {
		return nil, err
	}

	return r.(*Scheduler), nil
}

func (client Mikrotik) DeleteScheduler(name string) error {
	return client.Delete(&Scheduler{Name: name})
}
