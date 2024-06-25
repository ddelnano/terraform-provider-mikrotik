package client

import "reflect"

type (
	// FindByFieldWrapper changes the fields used to find the remote resource.
	FindByFieldWrapper struct {
		Resource
		field          string
		fieldValueFunc func() string
	}
)

var (
	_ Finder   = (*FindByFieldWrapper)(nil)
	_ Resource = (*FindByFieldWrapper)(nil)
)

func (fw FindByFieldWrapper) FindField() string {
	return fw.field
}

func (fw FindByFieldWrapper) FindFieldValue() string {
	return fw.fieldValueFunc()
}

// Create satisfies ResourceInstanceCreator interface and returns new object of the wrapped resource.
func (fw FindByFieldWrapper) Create() Resource {
	reflectNew := reflect.New(reflect.Indirect(reflect.ValueOf(fw.Resource)).Type())

	return reflectNew.Interface().(Resource)
}
