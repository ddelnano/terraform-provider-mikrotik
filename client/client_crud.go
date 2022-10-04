package client

import (
	"fmt"
	"log"
	"reflect"

	"github.com/go-routeros/routeros"
)

const (
	Add    Action = "add"
	Update Action = "update"
	Find   Action = "find"
	List   Action = "list"
	Delete Action = "delete"
)

type (
	// Action represents possible action on resource
	Action string

	// Resource interface defines a contract for abstract RouterOS resource
	Resource interface {
		// ActionToCommand trnaslates CRUD action to RouterOS command path
		ActionToCommand(Action) string

		// IDField reveals name of ID field to use in requests to MikroTik router
		// It is used in operations like Find
		IDField() string

		// ID returns value of the ID field
		ID() string

		// SetID updates a value of the ID field
		SetID(string)
	}

	// Adder defines contract for resources which require custom behaviour during resource creation
	Adder interface {
		// AfterAddHook is called right after the resource successfully added
		// This hook is mainly used to set resource's ID field based on reply from RouterOS
		AfterAddHook(r *routeros.Reply)
	}

	// Finder defines contract for resources which provide custom behaviour during resource retrieval
	Finder interface {
		// FindField retrieves a name of a field to use as key for resource retrieval
		FindField() string

		// FindFieldValue retrieves a value to use for resource retrieval
		FindFieldValue() string
	}

	// Deleter defines contract for resources which require custom behaviour during resource deletion
	Deleter interface {
		// DeleteField retrieves a name of a field which is used for resource deletion
		DeleteField() string

		// DeleteFieldValue retrieves a value for DeleteField field
		DeleteFieldValue() string
	}
)

// Add creates new resource on remote system
func (client Mikrotik) Add(d Resource) (interface{}, error) {
	c, err := client.getMikrotikClient()
	if err != nil {
		return nil, err
	}

	cmd := Marshal(d.ActionToCommand(Add), d)
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] creation response: `%v`", r)
	if adder, ok := d.(Adder); ok {
		adder.AfterAddHook(r)
	}

	return client.Find(d)
}

// Find retrieves resource from remote system
func (client Mikrotik) Find(d Resource) (interface{}, error) {
	findField := d.IDField()
	findFieldValue := d.ID()
	if finder, ok := d.(Finder); ok {
		findField = finder.FindField()
		findFieldValue = finder.FindFieldValue()
	}
	return client.findByField(d, findField, findFieldValue)
}

// Update updates existing resource on remote system
func (client Mikrotik) Update(resource Resource) (interface{}, error) {
	c, err := client.getMikrotikClient()
	if err != nil {
		return nil, err
	}

	cmd := Marshal(resource.ActionToCommand(Update), resource)
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	_, err = c.RunArgs(cmd)

	if err != nil {
		return nil, err
	}

	return client.Find(resource)
}

// Delete removes existing resource from remote system
func (client Mikrotik) Delete(d Resource) error {
	c, err := client.getMikrotikClient()
	if err != nil {
		return err
	}

	deleteField := d.IDField()
	deleteFieldValue := d.ID()
	if deleter, ok := d.(Deleter); ok {
		deleteField = deleter.DeleteField()
		deleteFieldValue = deleter.DeleteFieldValue()
	}
	cmd := []string{d.ActionToCommand(Delete), "=" + deleteField + "=" + deleteFieldValue}
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	_, err = c.RunArgs(cmd)
	if err != nil {
		return err
	}

	return nil
}

func (client Mikrotik) findByField(d Resource, field, value string) (interface{}, error) {
	cmd := []string{d.ActionToCommand(Find), "?" + field + "=" + value}
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)

	c, err := client.getMikrotikClient()
	if err != nil {
		return nil, err
	}
	r, err := c.RunArgs(cmd)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] find response: %v", r)

	targetStruct := client.newTargetStruct(d)
	targetStructInterface := targetStruct.Interface()
	err = Unmarshal(*r, targetStructInterface)
	if err != nil {
		return nil, err
	}
	if targetStructInterface.(Resource).ID() == "" {
		return nil, NewNotFound(fmt.Sprintf("resource with field `%s=%s` not found", field, value))
	}

	return targetStructInterface, nil
}

func (client Mikrotik) newTargetStruct(d interface{}) reflect.Value {
	return reflect.New(reflect.Indirect(reflect.ValueOf(d)).Type())
}
