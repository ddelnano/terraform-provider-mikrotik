package client

import (
	"github.com/go-routeros/routeros"
)

type Scheduler struct {
	Id        string `mikrotik:".id"`
	Name      string `mikrotik:"name"`
	OnEvent   string `mikrotik:"on-event"`
	StartDate string `mikrotik:"start-date"`
	StartTime string `mikrotik:"start-time"`
	Interval  int    `mikrotik:"interval,ttlToSeconds"`
}

var schedulerWrapper *resourceWrapper = &resourceWrapper{
	idField:       "name",
	idFieldDelete: "numbers",
	actionsMap: map[string]string{
		"add":    "/system/scheduler/add",
		"find":   "/system/scheduler/print",
		"update": "/system/scheduler/set",
		"delete": "/system/scheduler/remove",
	},
	targetStruct:          &Scheduler{},
	addIDExtractorFunc:    func(_ *routeros.Reply, resource interface{}) string { return resource.(*Scheduler).Name },
	recordIDExtractorFunc: func(r interface{}) string { return r.(*Scheduler).Name },
}

func (client Mikrotik) CreateScheduler(s *Scheduler) (*Scheduler, error) {
	r, err := schedulerWrapper.Add(s, client.getMikrotikClient)
	if err != nil {
		return nil, err
	}
	return r.(*Scheduler), nil
}

func (client Mikrotik) FindScheduler(name string) (*Scheduler, error) {
	r, err := schedulerWrapper.Find(name, client.getMikrotikClient)
	if err != nil {
		return nil, err
	}
	return r.(*Scheduler), nil

}

func (client Mikrotik) UpdateScheduler(s *Scheduler) (*Scheduler, error) {
	r, err := schedulerWrapper.Update(s, client.getMikrotikClient)
	if err != nil {
		return nil, err
	}
	return r.(*Scheduler), nil
}

func (client Mikrotik) DeleteScheduler(name string) error {
	return schedulerWrapper.Delete(name, client.getMikrotikClient)
}
