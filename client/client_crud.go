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
	// Action represents possible action on resource.
	Action string

	// Resource interface defines a contract for abstract RouterOS resource.
	Resource interface {
		// ActionToCommand translates CRUD action to RouterOS command path.
		ActionToCommand(Action) string

		// IDField reveals name of ID field to use in requests to MikroTik router.
		// It is used in operations like Find.
		IDField() string

		// ID returns value of the ID field.
		ID() string

		// SetID updates a value of the ID field.
		SetID(string)
	}

	// Adder defines contract for resources which require custom behaviour during resource creation.
	Adder interface {
		// AfterAddHook is called right after the resource successfully added.
		// This hook is mainly used to set resource's ID field based on reply from RouterOS.
		AfterAddHook(r *routeros.Reply)
	}

	// Finder defines contract for resources which provide custom behaviour during resource retrieval.
	Finder interface {
		// FindField retrieves a name of a field to use as key for resource retrieval.
		FindField() string

		// FindFieldValue retrieves a value to use for resource retrieval.
		FindFieldValue() string
	}

	// Deleter defines contract for resources which require custom behaviour during resource deletion.
	Deleter interface {
		// DeleteField retrieves a name of a field which is used for resource deletion.
		DeleteField() string

		// DeleteFieldValue retrieves a value for DeleteField field.
		DeleteFieldValue() string
	}

	// Normalizer is used to normalize response from RouterOS.
	// The main use-case is to populate fields which are empty in response but have default value,
	// for example `authoritative=yes` in `DHCPServer` resource is not returned by remote RouterOS instance.
	Normalizer interface {
		Normalize(r *routeros.Reply)
	}

	// ErrorHandler Defines contract to handle errors returned by RouterOS.
	// It can either return another error, or supress original error by returning nil.
	ErrorHandler interface {
		HandleError(error) error
	}

	// ResourceInstanceCreator interface defines methods to create new instance of a Resource.
	ResourceInstanceCreator interface {
		Create() Resource
	}
)

// Add creates new resource on remote system
func (client Mikrotik) Add(d Resource) (Resource, error) {
	c, err := client.getMikrotikClient()
	if err != nil {
		return nil, err
	}

	cmd := Marshal(d.ActionToCommand(Add), d)
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)
	if eh, ok := d.(ErrorHandler); ok {
		err = eh.HandleError(err)
	}
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
func (client Mikrotik) List(d Resource) ([]Resource, error) {
	cmd := []string{d.ActionToCommand(Find)}
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)

	c, err := client.getMikrotikClient()
	if err != nil {
		return nil, err
	}
	r, err := c.RunArgs(cmd)
	if eh, ok := d.(ErrorHandler); ok {
		err = eh.HandleError(err)
	}
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] find response: %v", r)

	targetStruct := client.newTargetStruct(d)
	targetSlicePtr := reflect.New(reflect.SliceOf(reflect.Indirect(targetStruct).Type()))
	targetSlice := reflect.Indirect(targetSlicePtr)
	err = Unmarshal(*r, targetSlicePtr.Interface())
	if err != nil {
		return nil, err
	}

	returnSlice := make([]Resource, 0)
	for i := 0; i < targetSlice.Len(); i++ {
		returnSlice = append(returnSlice, targetSlice.Index(i).Addr().Interface().(Resource))
	}

	return returnSlice, nil
}

// Find retrieves resource from remote system
func (client Mikrotik) Find(d Resource) (Resource, error) {
	findField := d.IDField()
	findFieldValue := d.ID()
	if finder, ok := d.(Finder); ok {
		findField = finder.FindField()
		findFieldValue = finder.FindFieldValue()
	}
	return client.findByField(d, findField, findFieldValue)
}

// Update updates existing resource on remote system
func (client Mikrotik) Update(resource Resource) (Resource, error) {
	c, err := client.getMikrotikClient()
	if err != nil {
		return nil, err
	}

	cmd := Marshal(resource.ActionToCommand(Update), resource)
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	_, err = c.RunArgs(cmd)
	if eh, ok := resource.(ErrorHandler); ok {
		err = eh.HandleError(err)
	}
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
	if rosErr, ok := err.(*routeros.DeviceError); ok {
		if rosErr.Sentence.Map["message"] == "no such item" {
			return NewNotFound(rosErr.Sentence.Map["message"])
		}
	}
	if eh, ok := d.(ErrorHandler); ok {
		err = eh.HandleError(err)
	}

	return err
}

func (client Mikrotik) findByField(d Resource, field, value string) (Resource, error) {
	cmd := []string{d.ActionToCommand(Find), "?" + field + "=" + value}
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)

	c, err := client.getMikrotikClient()
	if err != nil {
		return nil, err
	}
	r, err := c.RunArgs(cmd)
	if eh, ok := d.(ErrorHandler); ok {
		err = eh.HandleError(err)
	}
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] find response: %v", r)

	targetStruct := client.newTargetStruct(d)
	targetStructInterface := targetStruct.Interface()
	err = Unmarshal(*r, targetStructInterface)
	if eh, ok := d.(ErrorHandler); ok {
		err = eh.HandleError(err)
	}
	if err != nil {
		return nil, err
	}

	if n, ok := targetStructInterface.(Normalizer); ok {
		n.Normalize(r)
	}

	// assertion is not checked as we are creating the targetStruct from 'd' argument which satisfies Resource interface
	targetResource := targetStructInterface.(Resource)
	if targetResource.ID() == "" {
		return nil, NewNotFound(fmt.Sprintf("resource `%T` with field `%s=%s` not found", targetStruct, field, value))
	}

	return targetResource, nil
}

func (client Mikrotik) newTargetStruct(d interface{}) reflect.Value {
	if c, ok := d.(ResourceInstanceCreator); ok {
		return reflect.New(reflect.Indirect(reflect.ValueOf(c.Create())).Type())
	}

	return reflect.New(reflect.Indirect(reflect.ValueOf(d)).Type())
}
