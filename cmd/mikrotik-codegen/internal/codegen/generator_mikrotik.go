package codegen

import (
	"io"
)

func GenerateMikrotikResource(resourceName, commandBasePath string, w io.Writer) error {
	data := struct {
		CommandBasePath string
		ResourceName    string
	}{
		CommandBasePath: commandBasePath,
		ResourceName:    resourceName,
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
