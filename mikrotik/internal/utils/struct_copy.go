package utils

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// MikrotikStructToTerraformModel is a wrapper for copyStruct() to ensure proper src/dest typing
func MikrotikStructToTerraformModel(ctx context.Context, src client.Resource, dest interface{}) error {
	return copyStruct(ctx, src, dest)
}

// TerraformModelToMikrotikStruct is a wrapper for copyStruct() to ensure proper src/dest typing
func TerraformModelToMikrotikStruct(ctx context.Context, src interface{}, dest client.Resource) error {
	return copyStruct(ctx, src, dest)
}

// copyStruct copies fields of src struct to fields of dest struct.
//
// The fields matching is done based on field names (case insensitive).
// Having multiple fields with the same name but different case leads to unpredictable behavior.
//
// If dest struct has no field with particular name, it is skipped.
func copyStruct(ctx context.Context, src, dest interface{}) error {
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

		destField := reflectedDest.FieldByNameFunc(
			func(s string) bool {
				return strings.EqualFold(srcFieldType.Name, s)
			})
		destFieldType, found := reflectedDest.Type().FieldByNameFunc(
			func(s string) bool {
				return strings.EqualFold(srcFieldType.Name, s)
			})
		tflog.Debug(ctx, fmt.Sprintf("trying to copy struct field %q to %q", srcFieldType.Name, destFieldType.Name))
		if !destField.IsValid() || !found {
			// skip if dest struct does not have it (by name)
			tflog.Debug(ctx, "target field was not found")
			continue
		}
		if srcFieldType.PkgPath != "" || destFieldType.PkgPath != "" {
			// skip unexported fields
			tflog.Debug(ctx, "the source/target fields are unexported")
			continue
		}
		if !destField.CanSet() {
			// skip if dest field is not settable
			tflog.Debug(ctx, "target field is not settable")
			continue
		}

		switch kind := srcFieldType.Type.Kind(); kind {
		case reflect.Bool,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
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
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		tfValue = tftypes.Int64Value(int64(src.Uint()))
	case reflect.String:
		tfValue = tftypes.StringValue(src.String())
	case reflect.Bool:
		tfValue = tftypes.BoolValue(src.Bool())
	case reflect.Float32, reflect.Float64:
		tfValue = tftypes.Float64Value(src.Float())
	case reflect.Slice:
		var diags diag.Diagnostics
		var elements []interface{}
		for i := 0; i < src.Len(); i++ {
			elements = append(elements, src.Index(i).Interface())
		}
		var tfType attr.Type
		switch kind := src.Type().Elem().Kind(); kind {
		case reflect.Bool:
			tfType = tftypes.BoolType
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			tfType = tftypes.Int64Type
		case reflect.String:
			tfType = tftypes.StringType
		default:
			return fmt.Errorf("unsupported slice element type %q", kind)
		}
		var valueFromFunc func(t attr.Type, elements []interface{}) (attr.Value, diag.Diagnostics)

		switch dest.Interface().(type) {
		case tftypes.List:
			valueFromFunc = func(t attr.Type, elements []interface{}) (attr.Value, diag.Diagnostics) {
				return tftypes.ListValueFrom(context.TODO(), t, elements)
			}
		case tftypes.Set:
			valueFromFunc = func(t attr.Type, elements []interface{}) (attr.Value, diag.Diagnostics) {
				return tftypes.SetValueFrom(context.TODO(), t, elements)
			}
		default:
			return fmt.Errorf("unsupported destination Terraform type %v", reflect.TypeOf(dest).Name())
		}

		tfValue, diags = valueFromFunc(tfType, elements)

		if diags.HasError() {
			return fmt.Errorf("error creating Terraform type: %v", diags.Errors())
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
		case reflect.Uint:
			sliceType = reflect.TypeOf(uint(0))
		case reflect.Uint8:
			sliceType = reflect.TypeOf(uint8(0))
		case reflect.Uint16:
			sliceType = reflect.TypeOf(uint16(0))
		case reflect.Uint32:
			sliceType = reflect.TypeOf(uint32(0))
		case reflect.Uint64:
			sliceType = reflect.TypeOf(uint64(0))
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
	case tftypes.Set:
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
		case reflect.Uint:
			sliceType = reflect.TypeOf(uint(0))
		case reflect.Uint8:
			sliceType = reflect.TypeOf(uint8(0))
		case reflect.Uint16:
			sliceType = reflect.TypeOf(uint16(0))
		case reflect.Uint32:
			sliceType = reflect.TypeOf(uint32(0))
		case reflect.Uint64:
			sliceType = reflect.TypeOf(uint64(0))
		case reflect.String:
			sliceType = reflect.TypeOf("")
		default:
			return fmt.Errorf("unsupported list element types: %s -> []%s", src.Type().Name(), dest.Type().Elem().Kind())
		}
		targetPtr := reflect.New(reflect.SliceOf(sliceType))
		if len(f.Elements()) > 0 {
			diag = f.ElementsAs(context.TODO(), targetPtr.Interface(), false)
		}

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
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		dest.SetUint(src.Uint())
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
