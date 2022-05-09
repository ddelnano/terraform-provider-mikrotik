package codegen

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"reflect"
	"strings"
)

type (
	Struct struct {
		Name        string
		IDFieldName string
		Fields      []Field
	}

	Field struct {
		OriginalName string
		Name         string
		Tag          string
		Type         string
		Required     bool
		Optional     bool
		Computed     bool
	}
)

func Parse(fSet *token.FileSet, node ast.Node, startLine int, structName string) (*Struct, error) {
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

func parseStruct(structNode *ast.StructType) (*Struct, error) {
	result := &Struct{}

	for _, field := range structNode.Fields.List {
		tag := ""
		if field.Tag != nil {
			tag = field.Tag.Value
		}
		result.Fields = append(result.Fields,
			Field{
				Name: field.Names[0].Name,
				Tag:  tag,
				Type: fmt.Sprintf("%v", field.Type),
			},
		)
	}
	return result, nil
}

func parseStructUsingTags(structNode *ast.StructType) (*Struct, error) {
	result := &Struct{}

	for _, astField := range structNode.Fields.List {
		if astField.Tag == nil {
			continue
		}

		tag := reflect.StructTag(astField.Tag.Value)
		tagKey := "gen"
		tagValue, ok := tag.Lookup(tagKey)
		if !ok {
			continue
		}
		parts := strings.Split(tagValue, ",")
		name, opts := parts[0], parts[1:]
		var (
			optRequired = "required"
			optOptional = "optional"
			optComputed = "computed"
			optOmit     = "omit"
		)
		field := Field{
			OriginalName: astField.Names[0].Name,
			Name:         name,
			Tag:          tagValue,
			Type:         fmt.Sprintf("%v", astField.Type),
		}
		omit := false
		for _, o := range opts {
			switch {
			case o == optRequired:
				field.Required = true
			case o == optOptional:
				field.Optional = true
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
