package client

import (
	"fmt"
	"log"
	"strings"
)

type Script struct {
	Id                     string `mikrotik:".id"`
	Name                   string
	Owner                  string
	PolicyString           string `mikrotik:"policy"`
	DontRequirePermissions bool   `mikrotik:"dont-require-permissions"`
	Source                 string
}

func (s *Script) Policy() []string {
	return strings.Split(s.PolicyString, ",")
}

func (client Mikrotik) CreateScript(name, owner, source string, policies []string, dontReqPerms bool) (*Script, error) {
	c, err := client.getMikrotikClient()

	policiesString := strings.Join(policies, ",")
	nameArg := fmt.Sprintf("=name=%s", name)
	ownerArg := fmt.Sprintf("=owner=%s", owner)
	sourceArg := fmt.Sprintf("=source=%s", source)
	policyArg := fmt.Sprintf("=policy=%s", policiesString)
	dontReqPermsArg := fmt.Sprintf("=dont-require-permissions=%s", boolToMikrotikBool(dontReqPerms))
	cmd := []string{
		"/system/script/add",
		nameArg,
		ownerArg,
		sourceArg,
		policyArg,
		dontReqPermsArg,
	}
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)
	log.Printf("[DEBUG] /system/script/add returned %v", r)

	if err != nil {
		return nil, err
	}
	return client.FindScript(name)
}

func (client Mikrotik) UpdateScript(name, owner, source string, policy []string, dontReqPerms bool) (*Script, error) {
	c, err := client.getMikrotikClient()

	if err != nil {
		return nil, err
	}

	script, err := client.FindScript(name)

	if err != nil {
		return nil, err
	}

	policiesString := strings.Join(policy, ",")
	nameArg := fmt.Sprintf("=numbers=%s", script.Id)
	ownerArg := fmt.Sprintf("=owner=%s", owner)
	sourceArg := fmt.Sprintf("=source=%s", source)
	policyArg := fmt.Sprintf("=policy=%s", policiesString)
	dontReqPermsArg := fmt.Sprintf("=dont-require-permissions=%s", boolToMikrotikBool(dontReqPerms))
	cmd := []string{
		"/system/script/set",
		nameArg,
		ownerArg,
		sourceArg,
		policyArg,
		dontReqPermsArg,
	}
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	_, err = c.RunArgs(cmd)

	if err != nil {
		return nil, err
	}

	return client.FindScript(name)
}

func (client Mikrotik) DeleteScript(name string) error {
	c, err := client.getMikrotikClient()

	script, err := client.FindScript(name)

	if err != nil {
		return err
	}
	cmd := strings.Split(fmt.Sprintf("/system/script/remove =numbers=%s", script.Id), " ")
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)
	log.Printf("[DEBUG] Remove script from mikrotik api %v", r)

	return err
}

func (client Mikrotik) FindScript(name string) (*Script, error) {
	c, err := client.getMikrotikClient()
	cmd := strings.Split(fmt.Sprintf("/system/script/print ?name=%s", name), " ")
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)

	log.Printf("[DEBUG] Found script from mikrotik api %v", r)
	script := &Script{}
	err = Unmarshal(*r, script)

	if err != nil {
		return script, err
	}

	if script.Name == "" {
		return nil, NewNotFound(fmt.Sprintf("script `%s` not found", name))
	}

	return script, err
}

func boolToMikrotikBool(b bool) string {
	if b {
		return "yes"
	} else {
		return "no"
	}
}
