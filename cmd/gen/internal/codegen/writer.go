package codegen

import (
	"strings"
	"text/template"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	generatedNotice = "// Code generated by a script. DO NOT EDIT."

	resourceDefinitionTemplate = `
		package {{ .Package }}

		import (
			{{ range $import := .Imports -}}
				"{{ $import }}"
			{{ end }}
		)

		func resource{{ .ResourceName }}() *schema.Resource {
			return &schema.Resource{
				CreateContext: resource{{ .ResourceName }}Create,
				// ReadContext:   resource{{ .ResourceName }}Read,
				// UpdateContext: resource{{ .ResourceName }}Update,
				// DeleteContext: resource{{ .ResourceName }}Delete,
				Importer: &schema.ResourceImporter{
					StateContext: schema.ImportStatePassthroughContext,
				},

				Schema: map[string]*schema.Schema{
					{{ range .TerraformFields }}
					"{{ .Name }}": {
						Type:     schema.{{ .Type }},
						Required: {{ .Required }},
					},
					{{ end }}
				},
			}
		}

	func resource{{ .ResourceName }}Create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		record := dataTo{{ .ResourceName }}(d)

		c := m.(*client.Mikrotik)

		mikrotikRecord, err := c.Add{{ .ResourceName }}(record)
		if err != nil {
			return diag.FromErr(err)
		}

		return {{ .ResourceName | firstLower }}ToData(mikrotikRecord, d)
	}

	func dataTo{{ .ResourceName }}(d *schema.ResourceData) *client.{{ .ResourceName }} {
		record := new(client.{{ .ResourceName }})

		{{ range .Fields -}}
		record.{{ .Name }} = d.Get("{{ .Name | lowercase }}").({{ .Type }})
		{{ end }}
		return record
	}

	func {{ .ResourceName | firstLower }}ToData(record *client.{{ .ResourceName }}, d *schema.ResourceData) diag.Diagnostics {
		values := map[string]interface{}{
			{{ range .Fields -}}
			"{{ .Name | lowercase }}":    record.{{ .Name }},
			{{ end }}
		}

		d.SetId(record.{{ (index .Fields 0).Name }})

		var diags diag.Diagnostics

		for key, value := range values {
			if err := d.Set(key, value); err != nil {
				diags = append(diags, diag.Errorf("failed to set %s: %v", key, err)...)
			}
		}

		return diags
	}

	`
)

var (
	stringTypeToTerraformType = map[string]schema.ValueType{
		"string": schema.TypeString,
		"bool":   schema.TypeBool,
		"int":    schema.TypeInt,
	}
)

type (
	sourceWriter interface {
		Write([]byte) (int, error)
	}

	TerraformField struct {
		Name     string
		Type     schema.ValueType
		Required bool
	}

	templateData struct {
		Package         string
		Imports         []string
		ResourceName    string
		TerraformFields []TerraformField
		Fields          []Field
	}
)

func WriteSource(w sourceWriter, s Struct) error {
	if err := writeWrapper(w, []byte(generatedNotice)); err != nil {
		return err
	}
	fields := convertoToTerraformDefinition(s.Fields)
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

	if err := t.Execute(w,
		templateData{
			ResourceName:    s.Name,
			TerraformFields: fields,
			Fields:          s.Fields,
			Package:         "mikrotik",
			Imports: []string{
				"github.com/ddelnano/terraform-provider-mikrotik/client",
				"github.com/hashicorp/terraform-plugin-sdk/v2/diag",
				"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema",
			},
		}); err != nil {
		return err
	}

	return nil
}

func convertoToTerraformDefinition(fields []Field) []TerraformField {
	result := []TerraformField{}

	for _, f := range fields {
		result = append(result, TerraformField{
			Name:     f.Name,
			Type:     typeToTerraformType(f.Type),
			Required: true,
		})
	}

	return result
}

func typeToTerraformType(typ string) schema.ValueType {
	if t, ok := stringTypeToTerraformType[typ]; ok {
		return t
	}

	return schema.TypeInvalid
}

func writeWrapper(w sourceWriter, data []byte) error {
	_, err := w.Write(data)

	return err
}
