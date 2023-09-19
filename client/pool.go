package client

import (
	"github.com/go-routeros/routeros"
)

type Pool struct {
	Id       string `mikrotik:".id" codegen:"id,mikrotikID,terraformID"`
	Name     string `mikrotik:"name" codegen:"name,required"`
	Ranges   string `mikrotik:"ranges" codegen:"ranges,required"`
	NextPool string `mikrotik:"next-pool" codegen:"next_pool,optiona,computed"`
	Comment  string `mikrotik:"comment" codegen:"comment,optional,computed"`
}

var _ Resource = (*Pool)(nil)

func (b *Pool) ActionToCommand(a Action) string {
	return map[Action]string{
		Add:    "/ip/pool/add",
		Find:   "/ip/pool/print",
		Update: "/ip/pool/set",
		Delete: "/ip/pool/remove",
	}[a]
}

func (b *Pool) IDField() string {
	return ".id"
}

func (b *Pool) ID() string {
	return b.Id
}

func (b *Pool) SetID(id string) {
	b.Id = id
}

func (b *Pool) AfterAddHook(r *routeros.Reply) {
	b.Id = r.Done.Map["ret"]
}

// Typed wrappers
func (c Mikrotik) AddPool(r *Pool) (*Pool, error) {
	return r.processResourceErrorTuplePtr(c.Add(r))
}

func (c Mikrotik) UpdatePool(r *Pool) (*Pool, error) {
	return r.processResourceErrorTuplePtr(c.Update(r))
}

func (c Mikrotik) FindPool(id string) (*Pool, error) {
	return Pool{}.processResourceErrorTuplePtr(c.Find(&Pool{Id: id}))
}

func (c Mikrotik) FindPoolByName(name string) (*Pool, error) {
	return Pool{}.processResourceErrorTuplePtr(c.findByField(&Pool{}, "name", name))
}

func (c Mikrotik) DeletePool(id string) error {
	return c.Delete(&Pool{Id: id})
}

func (c Mikrotik) ListPools() ([]Pool, error) {
	res, err := c.List(&Pool{})
	if err != nil {
		return nil, err
	}
	returnSlice := make([]Pool, len(res))
	for i, v := range res {
		returnSlice[i] = *(v.(*Pool))
	}
	return returnSlice, nil
}

func (b Pool) processResourceErrorTuplePtr(r Resource, err error) (*Pool, error) {
	if err != nil {
		return nil, err
	}
	return r.(*Pool), nil
}
