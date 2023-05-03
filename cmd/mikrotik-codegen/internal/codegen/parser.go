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
	Struct struct {
		Name string
		// Original struct's field name to be used as Mikrotik resource ID
		MikrotikIDField string
		// Client's field to be used as Terraform resource ID
		TerraformIDField string
		DeleteField      string
		Fields           []Field
	}

	Field struct {
		OriginalName    string
		Name            string
		Type            string
		ElemType        string
		Required        bool
		Optional        bool
		DefaultValueStr string
		Computed        bool
	}
)

const (
	codegenTagKey = "codegen"

	optID         = "id"
	optMikrotikID = "mikrotikID"
	optDeleteID   = "deleteID"
	optRequired   = "required"
	optOptional   = "optional"
	optDefault    = "default="
	optType       = "type="
	optElemType   = "elemType="
	optComputed   = "computed"
	optOmit       = "omit"
)

// ParseFile parses a .go file with struct declaration.
//
// This functions searches for struct named `structName` or, in case of empty `structName`,
// finds first struct after the line where '//go:generate' is placed.
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

		field := Field{
			OriginalName: astField.Names[0].Name,
			Name:         name,
			Type:         fmt.Sprintf("%v", astField.Type),
		}
		omit := false
		for _, o := range opts {
			switch {
			case o == optID:
				if result.TerraformIDField != "" {
					return nil, fmt.Errorf("failed to set '%s' as Terraform ID field - it is already set to '%s'", field.OriginalName, result.TerraformIDField)
				}
				result.TerraformIDField = field.OriginalName
			case o == optMikrotikID:
				if result.MikrotikIDField != "" {
					return nil, fmt.Errorf("failed to set '%s' as Mikrotik ID field - it is already set to '%s'", field.OriginalName, result.MikrotikIDField)
				}
				result.MikrotikIDField = field.OriginalName
				// Mikrotik .id field should not appear in Terraform code
				omit = true
			case o == optDeleteID:
				if result.DeleteField != "" {
					return nil, fmt.Errorf("failed to set '%s' as delete ID field - it is already set to '%s'", field.OriginalName, result.DeleteField)
				}
				result.DeleteField = field.OriginalName
			case o == optRequired:
				field.Required = true
			case o == optOptional:
				field.Optional = true
			case strings.HasPrefix(o, optDefault):
				field.DefaultValueStr = strings.TrimPrefix(o, optDefault)
			case strings.HasPrefix(o, optType):
				field.Type = strings.TrimPrefix(o, optType)
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
		result.Fields = append(result.Fields, field)
	}

	return result, nil
}
