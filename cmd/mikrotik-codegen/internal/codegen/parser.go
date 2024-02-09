package codegen

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type (
	// Struct holds information about parsed struct.
	Struct struct {
		// Name is a of parsed struct.
		Name string

		// MikrotikIDField is a field name which holds MikroTik resource ID.
		MikrotikIDField string

		// TerraformIDField holds a field name which will be used as Terraform resource ID.
		TerraformIDField string

		// DeleteField holds a field name to use when deleting resource on MikroTik system.
		DeleteField string

		// Fields is a collection of field definitions in the parsed struct.
		Fields []*Field
	}

	// Field holds information about particular field in parsed struct.
	Field struct {
		// OriginalName is an original field name without chnages.
		OriginalName string

		// Name is a field name defined by struct tag.
		Name string

		// Required marks field as `required` in Terraform definition.
		Required bool

		// Optional marks field as `optional` in Terraform definition.
		Optional bool

		// Computed marks field as `computed` in Terraform definition.
		Computed bool

		// Type holds a field type.
		Type string

		// ElemType holds an element type if field type is List or Set
		ElemType string
	}
)

const (
	codegenTagKey = "codegen"

	optTerraformID = "terraformID"
	optMikrotikID  = "mikrotikID"
	optDeleteID    = "deleteID"
	optRequired    = "required"
	optOptional    = "optional"
	optComputed    = "computed"
	optElemType    = "elemType="
	optOmit        = "omit"
)

// ParseFile parses a .go file with struct declaration.
//
// This functions searches for struct definition `structName` and parses it.
// If `structName` is empty, function stops at first struct definition in the file right after `startLine`.
func ParseFile(filename string, startLine int, structName string) (*Struct, error) {
	_, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}

	fSet := token.NewFileSet()
	aFile, err := parser.ParseFile(fSet, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	if aFile == nil {
		return nil, errors.New("parsing of the file returned unexpected nil as *ast.File")
	}

	s, err := parse(fSet, aFile, startLine, structName)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func parse(fSet *token.FileSet, node ast.Node, startLine int, structName string) (*Struct, error) {
	structNode, foundName, err := findStruct(fSet, node, startLine, structName)
	if err != nil {
		return nil, err
	}

	parsedStruct, err := parseStructUsingTags(structNode)
	if err != nil {
		return nil, err
	}
	parsedStruct.Name = foundName

	return parsedStruct, nil
}

func findStruct(fSet *token.FileSet, node ast.Node, startLine int, structName string) (*ast.StructType, string, error) {
	var foundName string
	var structNode *ast.StructType

	ast.Inspect(node, func(n ast.Node) bool {
		if n == nil {
			return true
		}
		if n.Pos().IsValid() {
			pos := fSet.Position(n.Pos())
			if pos.Line < startLine {
				return true
			}
		}
		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}
		if typeSpec.Type == nil {
			return true
		}
		// if struct name is provided, ignore other structs on the way
		if structName != "" && typeSpec.Name.Name != structName {
			return true
		}

		foundName = typeSpec.Name.Name
		t, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			return true
		}

		structNode = t

		// stop after first struct is found
		return false
	})
	if foundName == "" {
		return nil, "", errors.New("struct not found")
	}
	return structNode, foundName, nil
}

func parseStructUsingTags(structNode *ast.StructType) (*Struct, error) {
	result := &Struct{}

	for _, astField := range structNode.Fields.List {
		if astField.Tag == nil {
			continue
		}

		// always unquote tag literal, otherwise it is treated as '`key:"options,here"`'
		unquoted, err := strconv.Unquote(astField.Tag.Value)
		if err != nil {
			return nil, err
		}
		tag := reflect.StructTag(unquoted)
		tagKey := codegenTagKey
		tagValue, ok := tag.Lookup(tagKey)
		if !ok {
			continue
		}
		parts := strings.Split(tagValue, ",")
		name, opts := parts[0], parts[1:]

		if name == "-" {
			continue
		}

		// determine the type of the field
		typeName := typeUnknown
		if exp, ok := astField.Type.(*ast.SelectorExpr); ok {
			// selector expression when type comes from another package, e.g. types.MikrotikList
			typeName = exp.Sel.Name
		}
		if exp, ok := astField.Type.(*ast.Ident); ok {
			// identifier, when it is a builtin type, e.g. "string"
			typeName = exp.Name
		}
		field := Field{
			OriginalName: astField.Names[0].Name,
			Name:         name,
			Type:         typeName,
		}
		omit := false
		for _, o := range opts {
			switch {
			case o == optTerraformID:
				if result.TerraformIDField != "" {
					return nil, fmt.Errorf("failed to set '%s' as Terraform ID field - it is already set to '%s'", field.OriginalName, result.TerraformIDField)
				}
				result.TerraformIDField = field.OriginalName
			case o == optMikrotikID:
				if result.MikrotikIDField != "" {
					return nil, fmt.Errorf("failed to set '%s' as Mikrotik ID field - it is already set to '%s'", field.OriginalName, result.MikrotikIDField)
				}
				result.MikrotikIDField = field.OriginalName
			case o == optDeleteID:
				if result.DeleteField != "" {
					return nil, fmt.Errorf("failed to set '%s' as delete ID field - it is already set to '%s'", field.OriginalName, result.DeleteField)
				}
				result.DeleteField = field.OriginalName
			case o == optRequired:
				field.Required = true
			case o == optOptional:
				field.Optional = true
			case strings.HasPrefix(o, optElemType):
				field.ElemType = strings.TrimPrefix(o, optElemType)
			case o == optComputed:
				field.Computed = true
			case o == optOmit:
				omit = true
			}
		}
		if omit {
			continue
		}
		if !(field.Computed || field.Required || field.Optional) {
			field.Optional = true
		}
		if field.OriginalName == result.MikrotikIDField {
			field.Computed = true
			field.Required = false
			field.Optional = false
		}

		result.Fields = append(result.Fields, &field)
	}

	if result.MikrotikIDField == "" {
		return nil, fmt.Errorf("MikroTik ID field is not set for any of the fields. Did you forget to mark one with '%s'?", optMikrotikID)
	}
	if result.TerraformIDField == "" {
		result.TerraformIDField = result.MikrotikIDField
	}

	return result, nil
}
