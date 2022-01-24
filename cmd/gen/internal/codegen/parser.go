package codegen

import (
	"fmt"
	"go/ast"
)

type (
	Struct struct {
		Name   string
		Fields []Field
	}

	Field struct {
		Name string
		Tag  string
		Type string
	}
)

func Parse(node ast.Node, structName string) (*Struct, error) {
	return findStruct(node, structName)
}

func findStruct(node ast.Node, structName string) (*Struct, error) {
	result := &Struct{}

	ast.Inspect(node, func(n ast.Node) bool {
		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}
		if typeSpec.Type == nil {
			return true
		}
		result.Name = typeSpec.Name.Name

		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			return true
		}

		for _, field := range structType.Fields.List {
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

		return false
	})
	return result, nil
}
