package client

import (
	"fmt"
	"log"
	"reflect"

	"github.com/go-routeros/routeros"
)

type (
	Resource interface {
		ActionCommand(string) string
		AddIDExtractionFunc(r *routeros.Reply) string
		IDField() string
		ID() string
		SetID(string)
		DeleteIDField() string
		SetDeleteID(string)
	}
)

func (client Mikrotik) Add(d Resource) (interface{}, error) {
	c, err := client.getMikrotikClient()
	if err != nil {
		return nil, err
	}

	cmd := Marshal(d.ActionCommand("add"), d)
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] creation response: `%v`", r)

	id := d.AddIDExtractionFunc(r)
	d.SetID(id)

	return client.Find(d)
}

func (client Mikrotik) Find(d Resource) (interface{}, error) {
	return client.findByField(d, d.IDField(), d.ID())
}

func (client Mikrotik) findByField(d Resource, field, value string) (interface{}, error) {
	cmd := []string{d.ActionCommand("find"), "?" + field + "=" + value}
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

func (client Mikrotik) Update(resource Resource) (interface{}, error) {
	c, err := client.getMikrotikClient()
	if err != nil {
		return nil, err
	}

	cmd := Marshal(resource.ActionCommand("update"), resource)
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	_, err = c.RunArgs(cmd)

	if err != nil {
		return nil, err
	}

	return client.Find(resource)

}

func (client Mikrotik) Delete(d Resource) error {
	c, err := client.getMikrotikClient()
	if err != nil {
		return err
	}

	cmd := []string{d.ActionCommand("delete"), "=" + d.DeleteIDField() + "=" + d.ID()}
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	_, err = c.RunArgs(cmd)
	if err != nil {
		return err
	}

	return nil
}

func (client Mikrotik) newTargetStruct(d interface{}) reflect.Value {
	return reflect.New(reflect.Indirect(reflect.ValueOf(d)).Type())
}

// func (client Mikrotik) newListOfTargetStructs(d interface{}) reflect.Value {
// 	elem := reflect.Indirect(reflect.ValueOf(d))
// 	listType := reflect.SliceOf(elem.Type())
// 	return reflect.New(listType)
// }
