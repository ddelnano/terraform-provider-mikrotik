package codegen

import (
	"io"
	"text/template"
)

func GenerateMikrotikResource(resourceName, commandBasePath string, w io.Writer) error {
	if err := writeWrapper(w, []byte(generatedNotice)); err != nil {
		return err
	}
	t := template.New("resource")
	if _, err := t.Parse(mikrotikResourceDefinitionTemplate); err != nil {
		return err
	}

	data := struct {
		CommandBasePath string
		ResourceName    string
	}{
		CommandBasePath: commandBasePath,
		ResourceName:    resourceName,
	}

	return t.Execute(w, data)
}
