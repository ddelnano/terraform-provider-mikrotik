package codegen

import (
	"errors"
	"io"
	"reflect"
	"strings"
	"text/template"

	"github.com/ddelnano/terraform-provider-mikrotik/client/types"
	"github.com/ddelnano/terraform-provider-mikrotik/cmd/mikrotik-codegen/internal/utils"
)

const (
	// List of declaration identifiers from ast package while parsing the source code
	AstVarTypeString = "string"
	AstVarTypeInt    = "int"
	AstVarTypeBool   = "bool"
)

var (
	// Handle custom types from the client package
	AstVarTypeMikrotikList     = reflect.TypeOf(types.MikrotikList{}).Name()
	AstVarTypeMikrotikIntList  = reflect.TypeOf(types.MikrotikIntList{}).Name()
	AstVarTypeMikrotikDuration = reflect.TypeOf(types.MikrotikDuration(0)).Name()

	terraformResourceImports = []string{
		"context",
		"github.com/ddelnano/terraform-provider-mikrotik/client",
		"github.com/hashicorp/terraform-plugin-framework/path",
		"github.com/hashicorp/terraform-plugin-framework/resource",
		"github.com/hashicorp/terraform-plugin-framework/resource/schema",
		"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier",
		"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier",
	}

	terraformResourceTestImports = []string{
		"fmt",
		"testing",
		"github.com/ddelnano/terraform-provider-mikrotik/client",
		"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource",
		"github.com/hashicorp/terraform-plugin-sdk/v2/terraform",
	}
)

type (
	// SourceWriteHookFunc defines a hook func to mutate source before writing to destination
	SourceWriteHookFunc func([]byte) ([]byte, error)

	sourceWriter interface {
		Write([]byte) (int, error)
	}

	terraformField struct {
		Name          string
		AttributeName string
		Required      bool
		Optional      bool
		Computed      bool
		MikrotikField *Field
		Type          Type
		ElemType      Type
	}

	templateData struct {
		Imports          []string
		ResourceName     string
		Fields           []*terraformField
		TerraformIDField *terraformField
		MikrotikIDField  *Field
		DeleteField      *terraformField
	}
)

// GenerateResource generates Terraform resource and writes it to specified output
func GenerateResource(s *Struct, w io.Writer) error {
	data, err := generateTemplateData(*s)
	if err != nil {
		return err
	}
	data.Imports = terraformResourceImports

	return generateCode(w,
		"resource",
		terraformResourceDefinitionTemplate,
		data,
	)
}

// GenerateResourceTest generates Terraform resource acceptance test and writes it to specified output
func GenerateResourceTest(s *Struct, w io.Writer) error {
	data, err := generateTemplateData(*s)
	if err != nil {
		return err
	}
	data.Imports = terraformResourceTestImports

	return generateCode(w,
		"resource_test",
		terraformResourceTestDefinitionTemplate,
		data,
	)
}

func generateTemplateData(s Struct) (templateData, error) {
	fields, err := convertToTerraformDefinition(s.Fields)
	if err != nil {
		return templateData{}, err
	}

	findTerraformFieldByName := func(fields []*terraformField, name string) *terraformField {
		for i := range fields {
			if fields[i].MikrotikField.OriginalName == name {
				return fields[i]
			}
		}
		return &terraformField{}
	}
	findMikrotikFieldByName := func(fields []*Field, name string) *Field {
		for i := range fields {
			if fields[i].OriginalName == name {
				return fields[i]
			}
		}
		return &Field{}
	}

	idField := findTerraformFieldByName(fields, s.MikrotikIDField)
	if idField.AttributeName == "" {
		return templateData{}, errors.New("The source struct does not provide information about ID field. Did you forget to mark one via 'id' tag?")
	}
	terraformIdField := findTerraformFieldByName(fields, s.TerraformIDField)
	if terraformIdField.AttributeName == "" {
		terraformIdField = idField
	}
	deleteField := findTerraformFieldByName(fields, s.DeleteField)
	if deleteField.AttributeName == "" {
		deleteField = idField
	}

	mikrotikIdField := findMikrotikFieldByName(s.Fields, s.MikrotikIDField)
	if mikrotikIdField.OriginalName == "" {
		mikrotikIdField = idField.MikrotikField
	}

	idField.Computed = true
	idField.Required = false
	idField.Optional = false

	return templateData{
		ResourceName:     s.Name,
		Fields:           fields,
		TerraformIDField: terraformIdField,
		MikrotikIDField:  mikrotikIdField,
		DeleteField:      deleteField,
	}, nil
}

func generateCode(w io.Writer, templateName, templateBody string, templateData interface{}) error {
	t := template.New(templateName)
	t.Funcs(template.FuncMap{
		"lowerCase":  strings.ToLower,
		"snakeCase":  utils.ToSnakeCase,
		"firstLower": utils.FirstLower,
		"sampleData": sampleData,
	})
	if _, err := t.Parse(templateBody); err != nil {
		return err
	}

	if err := writeWrapper(w, []byte(generatedNotice)); err != nil {
		return err
	}

	if err := t.Execute(w, templateData); err != nil {
		return err
	}

	return nil
}

func convertToTerraformDefinition(fields []*Field) ([]*terraformField, error) {
	result := []*terraformField{}

	for _, f := range fields {
		fieldType := typeToTerraformType(f.Type)
		elemType := UnknownType
		// currently, only list supports element typing
		if fieldType.Is(ListType) || fieldType.Is(SetType) {
			elemType = typeToTerraformType(f.ElemType)
		}
		result = append(result, &terraformField{
			Name:          f.OriginalName,
			AttributeName: f.Name,
			Type:          fieldType,
			ElemType:      elemType,
			Required:      f.Required,
			Optional:      f.Optional,
			Computed:      f.Computed,
			MikrotikField: f,
		})
	}

	return result, nil
}

func typeToTerraformType(typ string) Type {
	switch typ {
	case AstVarTypeBool:
		return BoolType
	case AstVarTypeInt, AstVarTypeMikrotikDuration:
		return Int64Type
	case AstVarTypeString:
		return StringType
	case AstVarTypeMikrotikList:
		return StringSliceType
	case AstVarTypeMikrotikIntList:
		return IntSliceType
	}

	return UnknownType
}

func writeWrapper(w sourceWriter, data []byte) error {
	_, err := w.Write(data)

	return err
}

// sampleData generates sample value for provided type.
func sampleData(typeName string) string {
	switch typeName {
	case typeString:
		return `"sample"`
	case typeList:
		return `[]`
	case typeSet:
		return `[]`
	case typeInt64:
		return "42"
	case typeBool:
		return "false"
	case typeStringSlice:
		return `["one", "two"]`
	case typeIntSlice:
		return `[1, 2, 3]`
	default:
		return `"` + typeUnknown + `"`
	}
}
