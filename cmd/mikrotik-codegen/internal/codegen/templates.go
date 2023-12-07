package codegen

const (
	generatedNotice = "// This code was generated. Review it carefully."

	terraformResourceDefinitionTemplate = `
package mikrotik

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
					{{.Type.Name | lowerCase}}planmodifier.UseStateForUnknown(),
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

	mikrotikResourceTestDefinitionTemplate = `
package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdd{{.ResourceName}}UpdateAndDelete(t *testing.T) {
	c := NewClient(GetConfigFromEnv())

	expectedResource := &{{.ResourceName}}{
	{{- range $field := .Fields }}
		{{- if and $field.Computed (not $field.Optional) }}{{continue}}{{end}}
		{{$field.Name}}: {{$field.Type.Name | sampleData}},
	{{- end }}
	}

	createdResource, err := c.Add{{.ResourceName}}(expectedResource)
	require.NoError(t, err)

	defer func(){
		id := createdResource.{{.TerraformIDField.Name}}
		err := c.Delete{{.ResourceName}}(id)
		if !assert.True(t, IsNotFoundError(err)) {
			assert.NoError(t, err)
		}
	}()

	expectedResource.Id = createdResource.Id

	foundResource, err := c.Find{{.ResourceName}}(expectedResource.{{.TerraformIDField.Name}})
	require.NoError(t, err)
	assert.Equal(t, expectedResource, foundResource)
{{ range $field := .Fields }}
	{{- if and $field.Computed (not $field.Optional) }}{{continue}}{{end}}
	expectedResource.{{$field.Name}} = expectedResource.{{$field.Name}} + {{$field.Type.Name | sampleData}}
{{- end }}

	updatedResource, err := c.Update{{.ResourceName}}(expectedResource)
	require.NoError(t, err)
	assert.Equal(t, expectedResource, updatedResource)

	// cleanup
	err = c.Delete{{.ResourceName}}(updatedResource.{{.TerraformIDField.Name}})
	assert.NoError(t, err)

	_, err = c.Find{{.ResourceName}}(expectedResource.{{.TerraformIDField.Name}})
	assert.Error(t, err)
}
`

	terraformResourceTestDefinitionTemplate = `
package mikrotik

import (
	{{ range $import := .Imports -}}
		"{{ $import }}"
	{{ end }}
)

{{ $resourceNameLower := .ResourceName | snakeCase }}
{{ $resourceType := .ResourceName | snakeCase | printf "mikrotik_%s" }}

func TestAcc{{.ResourceName}}_basic(t *testing.T) {
	resourceName := "{{$resourceType}}.testacc_{{$resourceNameLower}}"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheck{{.ResourceName}}Destroy,
		Steps: []resource.TestStep{
			{
				Config: ` + "`" + `
				resource "{{$resourceType}}" "testacc_{{$resourceNameLower}}" {
					{{ range $x := .Fields -}}
					{{if and $x.Computed (not $x.Optional) -}}{{continue}}{{end -}}
						{{ $x.AttributeName }} = {{$x.Type.Name | sampleData}}
					{{ end }}
				}
				` +
		"`" +
		`,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAcc{{.ResourceName}}Exists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("Not found: %s", resourceName)
					}
					return rs.Primary.Attributes["{{.TerraformIDField.AttributeName}}"], nil
				},
			},
		},
	})
}

func testAccCheck{{.ResourceName}}Destroy(s *terraform.State) error {
	c := client.NewClient(client.GetConfigFromEnv())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mikrotik_{{.ResourceName | snakeCase}}" {
			continue
		}

		remoteRecord, err := c.Find{{.ResourceName}}(rs.Primary.Attributes["{{.TerraformIDField.AttributeName}}"])

		if !client.IsNotFoundError(err) && err != nil {
			return err
		}

		if remoteRecord != nil {
			return fmt.Errorf("remote record (%s) still exists", remoteRecord.ID())
		}

	}
	return nil
}

func testAcc{{.ResourceName}}Exists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("%s does not exist in the statefile", resourceName)
		}

		c := client.NewClient(client.GetConfigFromEnv())
		record, err := c.Find{{.ResourceName}}(rs.Primary.Attributes["{{.TerraformIDField.AttributeName}}"])
		if err != nil {
			return fmt.Errorf("Unable to get remote record for %s: %v", resourceName, err)
		}

		if record == nil {
			return fmt.Errorf("Unable to get the remote record %s", resourceName)
		}

		return nil
	}
}

`
)
