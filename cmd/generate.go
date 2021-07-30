package main

import (
	"fmt"
	"io"
	"text/template"
)

var addResourceTemplate = `
func (client Mikrotik) Add{{ Name }}(l *{{ Type }}) (*{{ Type }}, error) {
	c, err := client.getMikrotikClient()

	if err != nil {
		return nil, err
	}

	cmd := Marshal("{{ MikrotikResourcePath }}/add", l)
	log.Printf("[INFO] Running the mikrotik command: %s", cmd)
	r, err := c.RunArgs(cmd)

	log.Printf("[DEBUG] {{ Name }} creation response: %v", r)

	if err != nil {
		return nil, err
	}

	id := r.Done.Map["ret"]

	return client.Find{{ Name }}(id)
}
`

var listResourceTemplate = `
func (client Mikrotik) List{{ Name }}() ([]{{ Name }}, error) {
	c, err := client.getMikrotikClient()

	if err != nil {
		return nil, err
	}
	cmd := []string{"{{ MikrotikResourcePath }}/print"}
	log.Printf("[INFO] Running the mikrotik command: %s", cmd)
	r, err := c.RunArgs(cmd)

	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] Found {{ Name }}s: %v", r)

	leases := []{{ Name }}{}

	err = Unmarshal(*r, &leases)

	if err != nil {
		return nil, err
	}

	return leases, nil
}

`

var findResourceTemplate = ``

var findResourceTemplateActual = `
func (client Mikrotik) Find{{ Name }}(id string) (*{{ Name }}, error) {
	c, err := client.getMikrotikClient()

	if err != nil {
		return nil, err
	}
	cmd := []string{"{{ MikrotikResourcePath }}/print", "?.id=" + id}
	log.Printf("[INFO] Running the mikrotik command: %s", cmd)
	r, err := c.RunArgs(cmd)

	log.Printf("[DEBUG] {{ Name }} response: %v", r)

	if err != nil {
		return nil, err
	}

	lease := {{ Name }}{}
	err = Unmarshal(*r, &lease)

	if err != nil {
		return nil, err
	}

	if lease.Id == "" {
		return nil, NewNotFound(fmt.Sprintf("{{ Name }} %s not found", id))
	}

	return &lease, nil
}
`

var updateResourceTemplate = `
func (client Mikrotik) Update{{ Name }}(l *{{ Name }} ) (*{{ Name }} , error) {
	c, err := client.getMikrotikClient()

	if err != nil {
		return nil, err
	}

	cmd := Marshal("{{ MikrotikResourcePath }}/set", l)
	log.Printf("[INFO] Running the mikrotik command: %s", cmd)
	_, err = c.RunArgs(cmd)

	if err != nil {
		return nil, err
	}

	return client.Find{{ Name }} (l.Id)
}
`

var deleteResourceTemplate = `
func (client Mikrotik) Delete{{ Name }}(id string) error {
	c, err := client.getMikrotikClient()

	if err != nil {
		return err
	}

	cmd := []string{"{{ MikrotikResourcePath }}/remove", "=.id=" + id}
	log.Printf("[INFO] Running the mikrotik command: %s", cmd)
	_, err = c.RunArgs(cmd)
	return err
}
`

func main() {
	templates := map[string]string{
		// "add":    addResourceTemplate,
		// "update": updateResourceTemplate,
		// "delete": deleteResourceTemplate,
		// "list":   listResourceTemplate,
		"find": findResourceTemplate,
	}
	t := template.New("")

	for tplType, tpl := range templates {
		_, err := t.New(tplType).Parse(tpl)
		if err != nil {
			fmt.Printf("failed to parse template: %v\n", err)
		}
	}

	err := t.ExecuteTemplate(io.Discard, "find", nil)
	if err != nil {
		fmt.Printf("failed to execute template: %v\n", err)
	}
}
