package codegen

import (
	"errors"
	"fmt"
	"go/ast"
)

type (
	Struct struct {
		Name        string
		IDFieldName string
		Fields      []Field
	}

	Field struct {
		Name string
		Tag  string
		Type string
	}
)

func Parse(node ast.Node, structName string) (*Struct, error) {
	structNode, err := findStruct(node, structName)
	if err != nil {
		return nil, err
	}

	parsedStruct, err := parseStruct(structNode)
	if err != nil {
		return nil, err
	}
	parsedStruct.Name = structName

	return parsedStruct, nil
}

func findStruct(node ast.Node, structName string) (*ast.StructType, error) {
	result := &Struct{}
	var structNode *ast.StructType

	ast.Inspect(node, func(n ast.Node) bool {
		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}
		if typeSpec.Type == nil {
			return true
		}
		if typeSpec.Name.Name != structName {
			return true
		}

		result.Name = typeSpec.Name.Name
		t, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			return true
		}
		structNode = t

		// stop after first struct is found
		return false
	})
	if result.Name == "" {
		return nil, errors.New("struct not found")
	}
	return structNode, nil
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
