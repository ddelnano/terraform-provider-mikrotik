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
		Name        string
		IDFieldName string
		Fields      []Field
	}

	Field struct {
		OriginalName string
		Name         string
		Type         string
		Required     bool
		Optional     bool
		Computed     bool
	}
)

const (
	optRequired = "required"
	optOptional = "optional"
	optComputed = "computed"
	optOmit     = "omit"
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
		tagKey := "gen"
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
