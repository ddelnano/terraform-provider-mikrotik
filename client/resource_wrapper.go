package client

import (
	"fmt"
	"log"
	"reflect"

	"github.com/go-routeros/routeros"
)

type (
	addIDExtractorFunc    func(r *routeros.Reply, resource interface{}) string
	recordIDExtractorFunc func(r interface{}) string
	mikrotikClientGetFunc func() (*routeros.Client, error)

	resourceWrapper struct {
		idField               string
		idFieldDelete         string
		actionsMap            map[string]string
		addIDExtractorFunc    addIDExtractorFunc
		recordIDExtractorFunc recordIDExtractorFunc
		targetStruct          interface{}
	}
)

func (rw *resourceWrapper) Add(resource interface{}, clientGetter mikrotikClientGetFunc) (interface{}, error) {
	c, err := clientGetter()
	if err != nil {
		return nil, err
	}

	cmd := Marshal(rw.actionsMap["add"], resource)
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] creation response: `%v`", r)

	id := rw.addIDExtractorFunc(r, resource)

	return rw.Find(id, clientGetter)
}

func (rw *resourceWrapper) Find(id string, clientGetter mikrotikClientGetFunc) (interface{}, error) {
	return rw.findByField(rw.idField, id, clientGetter)
}

func (rw *resourceWrapper) FindByField(field, value string, clientGetter mikrotikClientGetFunc) (interface{}, error) {
	return rw.findByField(field, value, clientGetter)
}

func (rw *resourceWrapper) findByField(field, value string, clientGetter mikrotikClientGetFunc) (interface{}, error) {
	cmd := []string{rw.actionsMap["find"], "?" + field + "=" + value}
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)

	c, err := clientGetter()
	if err != nil {
		return nil, err
	}
	r, err := c.RunArgs(cmd)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] find response: %v", r)

	targetStruct := rw.newTargetStruct()
	targetStructInterface := targetStruct.Interface()
	err = Unmarshal(*r, targetStructInterface)
	if err != nil {
		return nil, err
	}
	if rw.recordIDExtractorFunc(targetStructInterface) == "" {
		return nil, NewNotFound(fmt.Sprintf("resource with field `%s=%s` not found", field, value))
	}
	return targetStructInterface, nil
}

func (rw *resourceWrapper) List(clientGetter mikrotikClientGetFunc) (interface{}, error) {
	cmd := []string{rw.actionsMap["list"]}
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)

	c, err := clientGetter()
	if err != nil {
		return nil, err
	}
	r, err := c.RunArgs(cmd)

	log.Printf("[DEBUG] list response: %v", r)

	if err != nil {
		return nil, err
	}

	list := rw.newListOfTargetStructs()
	listInterface := list.Interface()
	err = Unmarshal(*r, listInterface)
	if err != nil {
		return nil, err
	}
	return reflect.Indirect(list).Interface(), nil
}

func (rw *resourceWrapper) Update(resource interface{}, clientGetter mikrotikClientGetFunc) (interface{}, error) {
	c, err := clientGetter()
	if err != nil {
		return nil, err
	}

	cmd := Marshal(rw.actionsMap["update"], resource)
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	_, err = c.RunArgs(cmd)

	if err != nil {
		return nil, err
	}

	return rw.Find(rw.recordIDExtractorFunc(resource), clientGetter)

}

func (rw *resourceWrapper) Delete(id string, clientGetter mikrotikClientGetFunc) error {
	c, err := clientGetter()
	if err != nil {
		return err
	}

	cmd := []string{rw.actionsMap["delete"], "=" + rw.idFieldDelete + "=" + id}
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	_, err = c.RunArgs(cmd)

	if err != nil {
		return err
	}

	return nil
}

func (rw *resourceWrapper) newTargetStruct() reflect.Value {
	return reflect.New(reflect.Indirect(reflect.ValueOf(rw.targetStruct)).Type())
}

func (rw *resourceWrapper) newListOfTargetStructs() reflect.Value {
	elem := reflect.Indirect(reflect.ValueOf(rw.targetStruct))
	listType := reflect.SliceOf(elem.Type())
	return reflect.New(listType)
}
