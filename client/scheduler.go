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
	return b.Name
}

// Typed wrappers
func (c Mikrotik) AddScheduler(r *Scheduler) (*Scheduler, error) {
	res, err := c.Add(r)
	if err != nil {
		return nil, err
	}

	return res.(*Scheduler), nil
}

func (c Mikrotik) UpdateScheduler(r *Scheduler) (*Scheduler, error) {
	res, err := c.Update(r)
	if err != nil {
		return nil, err
	}

	return res.(*Scheduler), nil
}

func (c Mikrotik) FindScheduler(name string) (*Scheduler, error) {
	res, err := c.Find(&Scheduler{Name: name})
	if err != nil {
		return nil, err
	}

	return res.(*Scheduler), nil
}

func (c Mikrotik) DeleteScheduler(name string) error {
	return c.Delete(&Scheduler{Name: name})
}
