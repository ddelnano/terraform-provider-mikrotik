package client

import (
	"fmt"
	"log"
	"reflect"

	"github.com/go-routeros/routeros"
)

type (
	addIDExtractorFunc    func(r *routeros.Reply) string
	recordIDExtractorFunc func(r interface{}) string
	mikrotikClientGetFunc func() (*routeros.Client, error)

	resourceWrapper struct {
		idField               string
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

	log.Printf("[DEBUG] ip address creation response: `%v`", r)

	if err != nil {
		return nil, err
	}

	id := rw.addIDExtractorFunc(r)

	return rw.Find(id, clientGetter)
}

func (rw *resourceWrapper) Find(id string, clientGetter mikrotikClientGetFunc) (interface{}, error) {
	cmd := []string{rw.actionsMap["find"], "?." + rw.idField + "=" + id}
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)

	c, err := clientGetter()
	if err != nil {
		return nil, err
	}
	r, err := c.RunArgs(cmd)

	log.Printf("[DEBUG] ip address response: %v", r)

	if err != nil {
		return nil, err
	}

	err = Unmarshal(*r, rw.targetStruct)
	if err != nil {
		return nil, err
	}
	if rw.recordIDExtractorFunc(rw.targetStruct) == "" {
		return nil, NewNotFound(fmt.Sprintf("resource `%s` not found", id))
	}
	return rw.targetStruct, nil
}

func (rw *resourceWrapper) List(clientGetter mikrotikClientGetFunc) (interface{}, error) {
	cmd := []string{rw.actionsMap["list"]}
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)

	c, err := clientGetter()
	if err != nil {
		return nil, err
	}
	r, err := c.RunArgs(cmd)

	log.Printf("[DEBUG] ip address response: %v", r)

	if err != nil {
		return nil, err
	}
	elem := reflect.Indirect(reflect.ValueOf(rw.targetStruct))
	listType := reflect.SliceOf(elem.Type())
	list := reflect.New(listType)
	err = Unmarshal(*r, list.Interface())
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

	cmd := []string{rw.actionsMap["delete"], "=." + rw.idField + "=" + id}
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	_, err = c.RunArgs(cmd)

	if err != nil {
		return err
	}

	return nil

}
