package utils

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
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
		_, ok := vDest.Type().FieldByName(srcFieldType.Name)
		if !ok {
			continue
		}
		if !destField.CanSet() {
			continue
		}

		tftypeCreateFunc := func(field reflect.Value) (interface{}, bool) {
			switch field.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				return tftypes.Int64Value(field.Int()), true
			case reflect.String:
				return tftypes.StringValue(field.String()), true
			case reflect.Bool:
				return tftypes.BoolValue(field.Bool()), true
			case reflect.Float32, reflect.Float64:
				return tftypes.Float64Value(field.Float()), true
			case reflect.Slice:
				var v interface{}
				var diag diag.Diagnostics
				var elements []interface{}
				for i := 0; i < field.Len(); i++ {
					elements = append(elements, field.Index(i).Interface())
				}
				switch field.Type().Elem().Kind() {
				case reflect.Bool:
					v, diag = tftypes.ListValueFrom(context.TODO(), tftypes.BoolType, elements)
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					v, diag = tftypes.ListValueFrom(context.TODO(), tftypes.Int64Type, elements)
				case reflect.String:
					v, diag = tftypes.ListValueFrom(context.TODO(), tftypes.StringType, elements)
				default:
					return nil, false
				}

				if diag.HasError() {
					return nil, false
				}
				return v, true
			}
			return nil, false
		}

		coreTypeToCoreType := func(src, dest reflect.Value) error {
			if src.Kind() != dest.Kind() {
				return errors.New("types mismatch")
			}

			switch src.Kind() {
			case reflect.Bool:
				dest.SetBool(src.Bool())
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				dest.SetInt(src.Int())
			case reflect.Float32, reflect.Float64:
				dest.SetFloat(src.Float())
			case reflect.String:
				dest.SetString(src.String())
			case reflect.Slice:
				slice := reflect.MakeSlice(dest.Type(), 0, 0)
				for i := 0; i < src.Len(); i++ {
					slice = reflect.Append(slice, src.Index(i))
				}
				dest.Set(slice)
			}

			return nil
		}

		switch srcFieldType.Type.Kind() {
		case reflect.Bool,
			reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8,
			reflect.Float32, reflect.Float64,
			reflect.String,
			reflect.Slice:
			// check this first as special case:
			// core type -> terraform type
			if attrValue, ok := destField.Interface().(attr.Value); ok {
				value, ok := tftypeCreateFunc(srcField)
				if !ok {
					return fmt.Errorf("unsupported mapping to Terraform type (field %q): %s -> %s",
						srcFieldType.Name,
						srcField.Kind(),
						attrValue.Type(context.TODO()).String(),
					)
				}
				destField.Set(reflect.ValueOf(value))
				break
			}

			// assume both are core types:
			// core type -> core type
			if err := coreTypeToCoreType(srcField, destField); err != nil {
				return err
			}
		}

		// source is terraform type and dest is core type

	}

	return nil
}
