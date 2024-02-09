package codegen

import (
	"io"
	"text/template"
	consoleinspected "github.com/ddelnano/terraform-provider-mikrotik/client/console-inspected"
	"github.com/ddelnano/terraform-provider-mikrotik/cmd/mikrotik-codegen/internal/utils"
)

func GenerateMikrotikResource(resourceName, commandBasePath string,
	consoleCommandDefinition consoleinspected.ConsoleItem,
	w io.Writer) error {
	if err := writeWrapper(w, []byte(generatedNotice)); err != nil {
		return err
	}
	t := template.New("resource")
	t.Funcs(template.FuncMap{
		"pascalCase": utils.PascalCase,
	})
	if _, err := t.Parse(mikrotikResourceDefinitionTemplate); err != nil {
		return err
	}

	fieldNames := make([]string, 0, len(consoleCommandDefinition.Arguments))
	for i := range consoleCommandDefinition.Arguments {
		fieldNames = append(fieldNames, consoleCommandDefinition.Arguments[i].Name)
	}

	data := struct {
		CommandBasePath string
		ResourceName    string
		FieldNames      []string
	}{
		CommandBasePath: commandBasePath,
		ResourceName:    resourceName,
		FieldNames:      fieldNames,
	}
	return generateCode(
		w,
		"resource",
		mikrotikResourceDefinitionTemplate,
		data,
	)
}

func GenerateMikrotikResourceTest(resourceName string, s *Struct, w io.Writer) error {
	data, err := generateTemplateData(*s)
	if err != nil {
		return err
	}

	return generateCode(
		w,
		"resource-test",
		mikrotikResourceTestDefinitionTemplate,
		data,
	)
}
