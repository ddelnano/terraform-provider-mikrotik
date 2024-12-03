package client

import (
	"github.com/ddelnano/terraform-provider-mikrotik/client/types"
	"github.com/go-routeros/routeros"
)

type Script struct {
	Id                     string             `mikrotik:".id" codegen:"id,deleteID"`
	Name                   string             `mikrotik:"name" codegen:"name,required,mikrotikID"`
	Owner                  string             `mikrotik:"owner,readonly" codegen:"owner,computed"`
	Policy                 types.MikrotikList `mikrotik:"policy" codegen:"policy,required"`
	DontRequirePermissions bool               `mikrotik:"dont-require-permissions" codegen:"dont_require_permissions"`
	Source                 string             `mikrotik:"source" codegen:"source,required"`
}

var _ Resource = (*Script)(nil)

func (b *Script) ActionToCommand(a Action) string {
	return map[Action]string{
		Add:    "/system/script/add",
		Find:   "/system/script/print",
		Update: "/system/script/set",
		Delete: "/system/script/remove",
	}[a]
}

func (b *Script) IDField() string {
	return ".id"
}

func (b *Script) ID() string {
	return b.Id
}

func (b *Script) SetID(id string) {
	b.Id = id
}

func (b *Script) AfterAddHook(r *routeros.Reply) {
	b.Id = r.Done.Map["ret"]
}

func (b *Script) FindField() string {
	return "name"
}

func (b *Script) FindFieldValue() string {
	return b.Name
}

func (b *Script) DeleteField() string {
	return "numbers"
}

func (b *Script) DeleteFieldValue() string {
	return b.Id
}

// Typed wrappers
func (c Mikrotik) AddScript(r *Script) (*Script, error) {
	res, err := c.Add(r)
	if err != nil {
		return nil, err
	}

	return res.(*Script), nil
}

func (c Mikrotik) UpdateScript(r *Script) (*Script, error) {
	res, err := c.Update(r)
	if err != nil {
		return nil, err
	}

	return res.(*Script), nil
}

func (c Mikrotik) FindScript(name string) (*Script, error) {
	res, err := c.Find(&Script{Name: name})
	if err != nil {
		return nil, err
	}

	return res.(*Script), nil
}

func (c Mikrotik) DeleteScript(id string) error {
	return c.Delete(&Script{Id: id})
}
