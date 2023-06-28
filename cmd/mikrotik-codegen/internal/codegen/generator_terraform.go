package codegen

import (
	"errors"
	"io"
	"strings"
	"text/template"

	"github.com/ddelnano/terraform-provider-mikrotik/cmd/mikrotik-codegen/internal/utils"
)

var (
	defaultImports = []string{
		"context",
		"github.com/ddelnano/terraform-provider-mikrotik/client",
		"github.com/hashicorp/terraform-plugin-framework/diag",
		"github.com/hashicorp/terraform-plugin-framework/path",
		"github.com/hashicorp/terraform-plugin-framework/resource",
		"github.com/hashicorp/terraform-plugin-framework/resource/schema",
		"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier",
		"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier",
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
		Package              string
		Imports              []string
		ResourceName         string
		Fields               []*terraformField
		TerraformIDField     *terraformField
		TerraformIDAttribute string
		MikrotikIDField      *Field
		DeleteField          *terraformField
	}
)

// GenerateResource generates Terraform resource and writes it to specified output
func GenerateResource(s *Struct, w io.Writer) error {
	if err := generateResource(w, *s); err != nil {
		return err
	}

	return nil
}

func generateResource(w sourceWriter, s Struct) error {
	if err := writeWrapper(w, []byte(generatedNotice)); err != nil {
		return err
	}
	fields, err := convertToTerraformDefinition(s.Fields)
	if err != nil {
		return err
	}

	t := template.New("resource")
	t.Funcs(template.FuncMap{
		"lowercase":  strings.ToLower,
		"snakeCase":  utils.ToSnakeCase,
		"firstLower": utils.FirstLower,
	})
	if _, err := t.Parse(resourceDefinitionTemplate); err != nil {
		return err
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
		return errors.New("The source struct does not provide information about ID field. Did you forget to mark one via 'id' tag?")
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
	if err := t.Execute(w,
		templateData{
			ResourceName:     s.Name,
			Fields:           fields,
			Package:          "mikrotik",
			TerraformIDField: terraformIdField,
			MikrotikIDField:  mikrotikIdField,
			DeleteField:      deleteField,
			Imports:          defaultImports,
		}); err != nil {
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
	case "slice":
		return ListType
	case "bool":
		return BoolType
	case "int":
		return Int64Type
	case "string":
		return StringType
	}

	return UnknownType
}

func writeWrapper(w sourceWriter, data []byte) error {
	_, err := w.Write(data)

	return err
}
