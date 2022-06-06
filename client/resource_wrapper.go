package client

import (
	"errors"
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
		mikrotikClientGetFunc mikrotikClientGetFunc
		actionsMap            map[string]string
		addIDExtractorFunc    addIDExtractorFunc
		recordIDExtractorFunc recordIDExtractorFunc
		targetStruct          interface{}
	}
)

func newResourceWrapper(mikrotikClientGetFunc mikrotikClientGetFunc, actionsMap map[string]string, idFieldExtractorFunc addIDExtractorFunc) (*resourceWrapper, error) {
	if mikrotikClientGetFunc == nil {
		return nil, errors.New("mikrotik client getter can not be nil")
	}

	return &resourceWrapper{
		mikrotikClientGetFunc: mikrotikClientGetFunc,
		actionsMap:            actionsMap,
		addIDExtractorFunc:    idFieldExtractorFunc,
	}, nil
}

func (rw *resourceWrapper) Add(resource interface{}) (interface{}, error) {
	c, err := rw.mikrotikClientGetFunc()
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

	return rw.Find(id)
}

func (rw *resourceWrapper) Find(id string) (interface{}, error) {
	cmd := []string{ipAddressWrapper.actionsMap["find"], "?." + rw.idField + "=" + id}
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)

	c, err := rw.mikrotikClientGetFunc()
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

func (rw *resourceWrapper) List() (interface{}, error) {
	cmd := []string{ipAddressWrapper.actionsMap["list"]}
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)

	c, err := rw.mikrotikClientGetFunc()
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

func (rw *resourceWrapper) Update(resource interface{}) (interface{}, error) {
	c, err := rw.mikrotikClientGetFunc()
	if err != nil {
		return nil, err
	}

	cmd := Marshal(rw.actionsMap["update"], resource)
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	_, err = c.RunArgs(cmd)

	if err != nil {
		return nil, err
	}

	return rw.Find(rw.recordIDExtractorFunc(resource))

}

func (rw *resourceWrapper) Delete(id string) error {
	c, err := rw.mikrotikClientGetFunc()
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

func (rw *resourceWrapper) WithMikrotikClientGetter(f mikrotikClientGetFunc) *resourceWrapper {
	rw.mikrotikClientGetFunc = f
	return rw
}
