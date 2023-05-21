package utils

import (
	"errors"
	"fmt"
	"reflect"
)

func CopyStruct(src, dest interface{}) error {
	if reflect.ValueOf(dest).Kind() != reflect.Pointer {
		return errors.New("dest must be a pointer")
	}

	vSrc := reflect.Indirect(reflect.ValueOf(src))
	vDest := reflect.ValueOf(dest).Elem()
	if vSrc.Kind() != reflect.Struct || vDest.Kind() != reflect.Struct {
		return fmt.Errorf("expected source and destination to be structs, got %v and %v", vSrc.Kind(), vDest.Kind())
	}
	for i := 0; i < vSrc.NumField(); i++ {
		srcField := vSrc.Field(i)
		srcFieldType := vSrc.Type().Field(i)
		destField := vDest.FieldByName(srcFieldType.Name)
		destFieldType, ok := vDest.Type().FieldByName(srcFieldType.Name)
		if !ok {
			continue
		}
		if srcFieldType.Type.Kind() != destFieldType.Type.Kind() {
			return errors.New("field types mismatch")
		}
		if !destField.CanSet() {
			continue
		}

		switch srcFieldType.Type.Kind() {
		case reflect.Bool:
			destField.SetBool(srcField.Bool())
		case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8:
			destField.SetInt(srcField.Int())
		case reflect.String:
			destField.SetString(srcField.String())
		case reflect.Slice:
			d := reflect.MakeSlice(srcField.Type(), 0, srcField.Len())
			for i := 0; i < srcField.Len(); i++ {
				d = reflect.Append(d, srcField.Index(i))
			}
			destField.Set(d)
		}
	}

	return nil
}
