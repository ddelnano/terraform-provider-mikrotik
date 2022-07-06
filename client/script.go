package client

import (
	"strings"

	"github.com/go-routeros/routeros"
)

type Script struct {
	Id                     string `mikrotik:".id"`
	Name                   string `mikrotik:"name"`
	Owner                  string `mikrotik:"owner"`
	PolicyString           string `mikrotik:"policy"`
	DontRequirePermissions bool   `mikrotik:"dont-require-permissions"`
	Source                 string `mikrotik:"source"`
}

var scriptWrapper *resourceWrapper = &resourceWrapper{
	idField:       "name",
	idFieldDelete: "numbers",
	actionsMap: map[string]string{
		"add":    "/system/script/add",
		"find":   "/system/script/print",
		"update": "/system/script/set",
		"delete": "/system/script/remove",
	},
	targetStruct:          &Script{},
	addIDExtractorFunc:    func(_ *routeros.Reply, resource interface{}) string { return resource.(*Script).Name },
	recordIDExtractorFunc: func(r interface{}) string { return r.(*Script).Name },
}

func (s *Script) Policy() []string {
	return strings.Split(s.PolicyString, ",")
}

func (client Mikrotik) CreateScript(name, owner, source string, policies []string, dontReqPerms bool) (*Script, error) {
	return client.AddScript(&Script{
		Name:                   name,
		Owner:                  owner,
		Source:                 source,
		PolicyString:           strings.Join(policies, ","),
		DontRequirePermissions: dontReqPerms,
	})
}

func (client Mikrotik) AddScript(s *Script) (*Script, error) {
	r, err := scriptWrapper.Add(s, client.getMikrotikClient)
	if err != nil {
		return nil, err
	}

	return r.(*Script), nil
}

func (client Mikrotik) FindScript(name string) (*Script, error) {
	r, err := scriptWrapper.Find(name, client.getMikrotikClient)
	if err != nil {
		return nil, err
	}

	return r.(*Script), nil
}

func (client Mikrotik) UpdateScript(name, owner, source string, policy []string, dontReqPerms bool) (*Script, error) {
	return client.updateScript(&Script{
		Name:                   name,
		Owner:                  owner,
		Source:                 source,
		PolicyString:           strings.Join(policy, ","),
		DontRequirePermissions: dontReqPerms,
	})
}

func (client Mikrotik) updateScript(s *Script) (*Script, error) {
	r, err := scriptWrapper.Update(s, client.getMikrotikClient)
	if err != nil {
		return nil, err
	}

	return r.(*Script), nil
}

func (client Mikrotik) DeleteScript(name string) error {
	return scriptWrapper.Delete(name, client.getMikrotikClient)
}
