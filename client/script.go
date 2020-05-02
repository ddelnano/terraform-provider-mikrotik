package client

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type Script struct {
	Id                     string
	Name                   string
	Owner                  string `mikrotik:-`
	Policy                 []string
	DontRequirePermissions bool
	Source                 string
}

func (client Mikrotik) CreateScript(name, owner, source string, policies []string, dontReqPerms bool) (Script, error) {
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
		return Script{}, err
	}
	return client.FindScript(name)
}

func (client Mikrotik) UpdateScript(name, owner, source string, policy []string, dontReqPerms bool) (Script, error) {
	c, err := client.getMikrotikClient()

	if err != nil {
		return Script{}, err
	}

	script, err := client.FindScript(name)

	if err != nil {
		return script, err
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
		return script, err
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

func (client Mikrotik) FindScript(name string) (Script, error) {
	c, err := client.getMikrotikClient()
	cmd := strings.Split(fmt.Sprintf("/system/script/print ?name=%s", name), " ")
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)

	log.Printf("[DEBUG] Found script from mikrotik api %v", r)

	if r.Re == nil {
		return Script{}, nil
	}
	if len(r.Re) > 1 && len(r.Re[0].List) > 1 {
		return Script{}, fmt.Errorf("Found more than one result for script with name %s", name)
	}
	script := Script{}
	for _, pair := range r.Re[0].List {
		if pair.Key == ".id" {
			script.Id = pair.Value
		}
		if pair.Key == "name" {
			script.Name = pair.Value
		}
		if pair.Key == "owner" {
			script.Owner = pair.Value
		}
		if pair.Key == "policy" {
			script.Policy = strings.Split(pair.Value, ",")
		}
		if pair.Key == "dont-require-permissions" {
			b, _ := strconv.ParseBool(pair.Value)
			if b {
				script.DontRequirePermissions = true
			} else {
				script.DontRequirePermissions = false
			}
		}
		if pair.Key == "source" {
			script.Source = pair.Value
		}
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
