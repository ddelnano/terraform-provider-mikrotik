package client

import (
	"github.com/go-routeros/routeros"
)

type Pool struct {
	Id      string `mikrotik:".id"`
	Name    string `mikrotik:"name"`
	Ranges  string `mikrotik:"ranges"`
	Comment string `mikrotik:"comment"`
}

var poolWrapper *resourceWrapper = &resourceWrapper{
	idField:       ".id",
	idFieldDelete: ".id",
	actionsMap: map[string]string{
		"add":    "/ip/pool/add",
		"find":   "/ip/pool/print",
		"list":   "/ip/pool/print",
		"update": "/ip/pool/set",
		"delete": "/ip/pool/remove",
	},
	targetStruct:          &Pool{},
	addIDExtractorFunc:    func(r *routeros.Reply, _ interface{}) string { return r.Done.Map["ret"] },
	recordIDExtractorFunc: func(r interface{}) string { return r.(*Pool).Id },
}

func (client Mikrotik) AddPool(p *Pool) (*Pool, error) {
	r, err := poolWrapper.Add(p, client.getMikrotikClient)
	if err != nil {
		return nil, err
	}

	return r.(*Pool), nil
}

func (client Mikrotik) ListPools() ([]Pool, error) {
	r, err := poolWrapper.List(client.getMikrotikClient)
	if err != nil {
		return nil, err
	}

	return r.([]Pool), nil
}

func (client Mikrotik) FindPool(id string) (*Pool, error) {
	r, err := poolWrapper.Find(id, client.getMikrotikClient)
	if err != nil {
		return nil, err
	}

	return r.(*Pool), nil
}

func (client Mikrotik) FindPoolByName(name string) (*Pool, error) {
	r, err := poolWrapper.findByField("name", name, client.getMikrotikClient)
	if err != nil {
		return nil, err
	}

	return r.(*Pool), nil
}

func (client Mikrotik) UpdatePool(p *Pool) (*Pool, error) {
	r, err := poolWrapper.Update(p, client.getMikrotikClient)
	if err != nil {
		return nil, err
	}

	return r.(*Pool), nil
}

func (client Mikrotik) DeletePool(id string) error {
	return poolWrapper.Delete(id, client.getMikrotikClient)
}
