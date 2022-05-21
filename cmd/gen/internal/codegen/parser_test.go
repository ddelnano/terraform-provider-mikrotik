package codegen

import (
	"go/parser"
	"go/token"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	cases := []struct {
		name       string
		source     []byte
		structName string
		startLine  int
		expected   *Struct
	}{
		{
			name: "struct name provided",
			source: []byte(`
package testpackage

type DnsRecord struct {
	Name 			   string` + " `gen:\"name,required\"`" + `
	GeneratedNumber	   string` + " `gen:\"internal_id,computed\"`" + `
	Enabled 		   bool` + " `gen:\"enabled,optional\"`" + `
	ExplicitlyOmitted  bool` + " `gen:\"-,omit\"`" + `
}
			`),

			expected: &Struct{
				Name: "DnsRecord",
				Fields: []Field{
					{
						OriginalName: "Name",
						Name:         "name",
						Type:         "string",
						Required:     true,
					},
					{
						OriginalName: "GeneratedNumber",
						Name:         "internal_id",
						Type:         "string",
						Computed:     true,
					},
					{
						OriginalName: "Enabled",
						Name:         "enabled",
						Type:         "bool",
						Optional:     true,
					},
				},
			},
		},
		{
			name: "id field is parsed",
			source: []byte(`
package testpackage

type DnsRecord struct {
	Name 			   string` + " `gen:\"name,id,required\"`" + `
	GeneratedNumber	   string` + " `gen:\"internal_id,computed\"`" + `
	Enabled 		   bool` + " `gen:\"enabled,optional\"`" + `
	ExplicitlyOmitted  bool` + " `gen:\"-,omit\"`" + `
}
			`),

			expected: &Struct{
				Name:        "DnsRecord",
				IDFieldName: "Name",
				Fields: []Field{
					{
						OriginalName: "Name",
						Name:         "name",
						Type:         "string",
						Required:     true,
					},
					{
						OriginalName: "GeneratedNumber",
						Name:         "internal_id",
						Type:         "string",
						Computed:     true,
					},
					{
						OriginalName: "Enabled",
						Name:         "enabled",
						Type:         "bool",
						Optional:     true,
					},
				},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fSet := token.NewFileSet()
			node, err := parser.ParseFile(fSet, "", tc.source, parser.ParseComments)
			if err != nil {
				t.Error(err)
			}
			result, err := parse(fSet, node, tc.startLine, tc.structName)
			if err != nil {
				t.Error(err)
			}
			if !reflect.DeepEqual(tc.expected, result) {
				t.Errorf(`
				objects differ:
					wanted:
						%+#v
					got:
						%+#v
					`, tc.expected, result)
			}
		})
	}
}
