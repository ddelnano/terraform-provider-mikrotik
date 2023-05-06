package codegen

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"text/template"
)

var (
	// we should find a better (typed) way to represent this mapping
	stringTypeToTerraformType = map[string]string{
		"slice":  "List",
		"string": "String",
		"bool":   "Bool",
		"int":    "Int64",
	}

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
		Type          string
		ElemType      string
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
func GenerateResource(s *Struct, w io.Writer, beforeWriteHooks ...SourceWriteHookFunc) error {
	var result []byte
	var buf bytes.Buffer
	var err error

	if err := generateResource(&buf, *s); err != nil {
		return err
	}
	result = buf.Bytes()
	for _, h := range beforeWriteHooks {
		result, err = h(result)
		if err != nil {
			return err
		}
	}

	_, err = w.Write(result)
	if err != nil {
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
		"lowercase": strings.ToLower,
		"firstLower": func(s string) string {
			if len(s) < 1 {
				return s
			}
			if len(s) == 1 {
				return strings.ToLower(s)
			}

			return strings.ToLower(s[:1]) + s[1:]
		},
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
		elemType := "String"
		// currently, only list supports element typing
		if fieldType == "List" || fieldType == "Set" {
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

func typeToTerraformType(typ string) string {
	if t, ok := stringTypeToTerraformType[typ]; ok {
		return t
	}

	return stringTypeToTerraformType["string"]
}

func writeWrapper(w sourceWriter, data []byte) error {
	_, err := w.Write(data)

	return err
}
