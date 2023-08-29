package codegen

const (
	generatedNotice = "// This code was generated. Review it carefully."

	resourceDefinitionTemplate = `
package {{ .Package }}

import (
	{{ range $import := .Imports -}}
		"{{ $import }}"
	{{ end }}
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"

)
{{ $resourceStructName := .ResourceName | firstLower}}
type {{$resourceStructName}} struct {
	client *client.Mikrotik
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &{{$resourceStructName}}{}
	_ resource.ResourceWithConfigure   = &{{$resourceStructName}}{}
	_ resource.ResourceWithImportState = &{{$resourceStructName}}{}
)

// New{{.ResourceName}}Resource is a helper function to simplify the provider implementation.
func New{{.ResourceName}}Resource() resource.Resource {
	return &{{$resourceStructName}}{}
}

func (r *{{$resourceStructName}}) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*client.Mikrotik)
}

// Metadata returns the resource type name.
func (r *{{$resourceStructName}}) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_{{.ResourceName | snakeCase}}"
}

// Schema defines the schema for the resource.
func (s *{{$resourceStructName}}) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates a MikroTik {{.ResourceName}}.",
		Attributes: map[string]schema.Attribute{
			{{range .Fields -}}
			"{{.AttributeName}}": schema.{{.Type.Name}}Attribute{
				Required: {{.Required}},
				Optional: {{.Optional}},
				Computed: {{.Computed}},
				{{if .Computed -}}
				PlanModifiers: []planmodifier.{{.Type.Name}}{
					{{.Type.Name | lowercase}}planmodifier.UseStateForUnknown(),
				},
				{{- end}}
				Description: "",
			},
			{{end}}
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *{{$resourceStructName}}) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var terraformModel {{$resourceStructName}}Model
	var mikrotikModel client.{{.ResourceName}}
	GenericCreateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Read refreshes the Terraform state with the latest data.
func (r *{{$resourceStructName}}) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var terraformModel {{$resourceStructName}}Model
	var mikrotikModel client.{{.ResourceName}}
	GenericReadResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *{{$resourceStructName}}) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var terraformModel {{$resourceStructName}}Model
	var mikrotikModel client.{{.ResourceName}}
	GenericUpdateResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *{{$resourceStructName}}) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var terraformModel {{$resourceStructName}}Model
	var mikrotikModel client.{{.ResourceName}}
	GenericDeleteResource(&terraformModel, &mikrotikModel, r.client)(ctx, req, resp)
}

func (r *{{$resourceStructName}}) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("{{.TerraformIDField.AttributeName}}"), req, resp)
}

type {{$resourceStructName}}Model struct {
	{{range .Fields -}}
	{{.Name}}        tftypes.{{.Type.Name}} ` + "`" + `tfsdk:"{{.AttributeName}}"` + "`" + `
	{{end}}
}
`

	mikrotikResourceDefinitionTemplate = `
package client

import (
	"github.com/ddelnano/terraform-provider-mikrotik/client/internal/types"
	"github.com/go-routeros/routeros"
)

// {{.ResourceName}} defines resource
type {{.ResourceName}} struct {
	Id string ` + "`" + `mikrotik:".id"` + "`" + `
}

var _ Resource = (*{{.ResourceName}})(nil)

func (b *{{.ResourceName}}) ActionToCommand(a Action) string {
	return map[Action]string{
		Add:    "{{.CommandBasePath}}/add",
		Find:   "{{.CommandBasePath}}/print",
		Update: "{{.CommandBasePath}}/set",
		Delete: "{{.CommandBasePath}}/remove",
	}[a]
}

func (b *{{.ResourceName}}) IDField() string {
	return ".id"
}

func (b *{{.ResourceName}}) ID() string {
	return b.Id
}

func (b *{{.ResourceName}}) SetID(id string) {
	b.Id = id
}

// Uncomment extra methods to satisfy more interfaces

// Adder
// func (b *{{.ResourceName}}) AfterAddHook(r *routeros.Reply) {
// 	b.Id = r.Done.Map["ret"]
// }

// Finder
// func (b *{{.ResourceName}}) FindField() string {
// 	return "name"
// }

// func (b *{{.ResourceName}}) FindFieldValue() string {
// 	return b.Name
// }

// Deleter
// func (b *{{.ResourceName}}) DeleteField() string {
// 	return "numbers"
// }

// func (b *{{.ResourceName}}) DeleteFieldValue() string {
// 	return b.Id
// }


// Typed wrappers
func (c Mikrotik) Add{{.ResourceName}}(r *{{.ResourceName}}) (*{{.ResourceName}}, error) {
	res, err := c.Add(r)
	if err != nil {
		return nil, err
	}

	return res.(*{{.ResourceName}}), nil
}

func (c Mikrotik) Update{{.ResourceName}}(r *{{.ResourceName}}) (*{{.ResourceName}}, error) {
	res, err := c.Update(r)
	if err != nil {
		return nil, err
	}

	return res.(*{{.ResourceName}}), nil
}

func (c Mikrotik) Find{{.ResourceName}}(id string) (*{{.ResourceName}}, error) {
	res, err := c.Find(&{{.ResourceName}}{Id: id})
	if err != nil {
		return nil, err
	}

	return res.(*{{.ResourceName}}), nil
}

func (c Mikrotik) List{{.ResourceName}}() ([]{{.ResourceName}}, error) {
	res, err := c.List(&{{.ResourceName}}{})
	if err != nil {
		return nil, err
	}
	returnSlice := make([]{{.ResourceName}}, len(res))
	for i, v := range res {
		returnSlice[i] = *(v.(*{{.ResourceName}}))
	}

	return returnSlice, nil
}


func (c Mikrotik) Delete{{.ResourceName}}(id string) error {
	return c.Delete(&{{.ResourceName}}{Id: id})
}

`
)
