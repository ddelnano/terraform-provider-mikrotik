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

// CopyStruct copies fields of src struct to fields of dest struct.
//
// The fields matching is done based on field names.
// If dest struct has no field with particular name, it is skipped.
func CopyStruct(src, dest interface{}) error {
	if reflect.ValueOf(dest).Kind() != reflect.Pointer {
		return errors.New("destination must be a pointer")
	}

	reflectedSrc := reflect.Indirect(reflect.ValueOf(src))
	reflectedDest := reflect.Indirect(reflect.ValueOf(dest))
	if reflectedSrc.Kind() != reflect.Struct || reflectedDest.Kind() != reflect.Struct {
		return fmt.Errorf("source and destination must be structs, got %v and %v", reflectedSrc.Kind(), reflectedDest.Kind())
	}

	for i := 0; i < reflectedSrc.NumField(); i++ {
		srcField := reflectedSrc.Field(i)
		srcFieldType := reflectedSrc.Type().Field(i)
		destField := reflectedDest.FieldByName(srcFieldType.Name)

		_, ok := reflectedDest.Type().FieldByName(srcFieldType.Name)
		if !ok {
			// skip if dest struct does not have it (by name)
			continue
		}
		if !destField.CanSet() {
			// skip if dest field is not settable
			continue
		}

		switch kind := srcFieldType.Type.Kind(); kind {
		case reflect.Bool,
			reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8,
			reflect.Float32, reflect.Float64,
			reflect.String,
			reflect.Slice:
			// core type -> terraform type
			// check if dest field is one of the Terraform types
			if _, ok := destField.Interface().(attr.Value); ok {
				if err := coreTypeToTerraformType(srcField, destField); err != nil {
					return err
				}
				break
			}

			// core type -> core type
			if err := coreTypeToCoreType(srcField, destField); err != nil {
				return err
			}
		case reflect.Struct:
			// source is terraform type and dest is core type
			if _, ok := srcField.Interface().(attr.Value); ok {
				if err := terraformTypeToCoreType(srcField, destField); err != nil {
					return err
				}
				break
			}
			return errors.New("unsupported source field of type 'struct'")
		default:
			return fmt.Errorf("unsupported source field type %q", kind)
		}
	}

	return nil
}

func coreTypeToTerraformType(src, dest reflect.Value) error {
	var tfValue attr.Value
	switch src.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		tfValue = tftypes.Int64Value(src.Int())
	case reflect.String:
		tfValue = tftypes.StringValue(src.String())
	case reflect.Bool:
		tfValue = tftypes.BoolValue(src.Bool())
	case reflect.Float32, reflect.Float64:
		tfValue = tftypes.Float64Value(src.Float())
	case reflect.Slice:
		var diag diag.Diagnostics
		var elements []interface{}
		for i := 0; i < src.Len(); i++ {
			elements = append(elements, src.Index(i).Interface())
		}
		switch kind := src.Type().Elem().Kind(); kind {
		case reflect.Bool:
			tfValue, diag = tftypes.ListValueFrom(context.TODO(), tftypes.BoolType, elements)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			tfValue, diag = tftypes.ListValueFrom(context.TODO(), tftypes.Int64Type, elements)
		case reflect.String:
			tfValue, diag = tftypes.ListValueFrom(context.TODO(), tftypes.StringType, elements)
		default:
			return fmt.Errorf("unsupported slice element type %q", kind)
		}

		if diag.HasError() {
			return fmt.Errorf("error creating Terraform type: %v", diag.Errors())
		}
	}

	dest.Set(reflect.ValueOf(tfValue))

	return nil
}

func terraformTypeToCoreType(src, dest reflect.Value) error {
	switch f := src.Interface().(type) {
	case tftypes.Int64:
		dest.SetInt(f.ValueInt64())
	case tftypes.String:
		dest.SetString(f.ValueString())
	case tftypes.Bool:
		dest.SetBool(f.ValueBool())
	case tftypes.List:
		var diag diag.Diagnostics
		var sliceType reflect.Type

		switch dest.Type().Elem().Kind() {
		case reflect.Bool:
			sliceType = reflect.TypeOf(true)
		case reflect.Int:
			sliceType = reflect.TypeOf(int(0))
		case reflect.Int8:
			sliceType = reflect.TypeOf(int8(0))
		case reflect.Int16:
			sliceType = reflect.TypeOf(int16(0))
		case reflect.Int32:
			sliceType = reflect.TypeOf(int32(0))
		case reflect.Int64:
			sliceType = reflect.TypeOf(int64(0))
		case reflect.String:
			sliceType = reflect.TypeOf("")
		default:
			return fmt.Errorf("unsupported list element types: %s -> []%s", src.Type().Name(), dest.Type().Elem().Kind())
		}
		targetPtr := reflect.New(reflect.SliceOf(sliceType))
		diag = f.ElementsAs(context.TODO(), targetPtr.Interface(), false)

		if diag.HasError() {
			return fmt.Errorf("%s", diag.Errors())
		}

		dest.Set(targetPtr.Elem())

		return nil
	default:
		return fmt.Errorf("unsupported field type assignment: %s -> %s", src.Type().Name(), dest.Kind())
	}

	return nil
}

func coreTypeToCoreType(src, dest reflect.Value) error {
	if src.Kind() != dest.Kind() {
		return fmt.Errorf("cannot assign %s to %s", src.Kind(), dest.Kind())
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
