package codegen

import (
	"errors"
	"go/parser"
	"go/token"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	cases := []struct {
		name          string
		source        []byte
		structName    string
		startLine     int
		expected      *Struct
		expectedError error
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
			name: "terraform and mikrotik id fields are parsed",
			source: []byte(`
package testpackage

type DnsRecord struct {
	ID	 			   string` + " `gen:\"-,mikrotikID\"`" + `
	Name 			   string` + " `gen:\"name,id,required\"`" + `
	GeneratedNumber	   string` + " `gen:\"internal_id,computed\"`" + `
	Enabled 		   bool` + " `gen:\"enabled,optional\"`" + `
	ExplicitlyOmitted  bool` + " `gen:\"-,omit\"`" + `
}
			`),

			expected: &Struct{
				Name:             "DnsRecord",
				TerraformIDField: "Name",
				MikrotikIDField:  "ID",
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
			name: "deleteID field is parsed",
			source: []byte(`
package testpackage

type DnsRecord struct {
	ID	 			   string` + " `gen:\"-,mikrotikID\"`" + `
	Name 			   string` + " `gen:\"name,id,deleteID,required\"`" + `
	GeneratedNumber	   string` + " `gen:\"internal_id,computed\"`" + `
	Enabled 		   bool` + " `gen:\"enabled,optional\"`" + `
	ExplicitlyOmitted  bool` + " `gen:\"-,omit\"`" + `
}
			`),

			expected: &Struct{
				Name:             "DnsRecord",
				TerraformIDField: "Name",
				MikrotikIDField:  "ID",
				DeleteField:      "Name",
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
			name: "default values are parsed",
			source: []byte(`
package testpackage

type DnsRecord struct {
	ID	 			   string` + " `gen:\"-,mikrotikID\"`" + `
	Name 			   string` + " `gen:\"name,id,deleteID,required\"`" + `
	GeneratedNumber	   int` + " `gen:\"internal_id,optional,default=10\"`" + `
	Enabled 		   bool` + " `gen:\"enabled,optional,default=true\"`" + `
	ExplicitlyOmitted  bool` + " `gen:\"-,omit\"`" + `
}
			`),

			expected: &Struct{
				Name:             "DnsRecord",
				TerraformIDField: "Name",
				MikrotikIDField:  "ID",
				DeleteField:      "Name",
				Fields: []Field{
					{
						OriginalName: "Name",
						Name:         "name",
						Type:         "string",
						Required:     true,
					},
					{
						OriginalName:    "GeneratedNumber",
						Name:            "internal_id",
						Type:            "int",
						Optional:        true,
						DefaultValueStr: "10",
					},
					{
						OriginalName:    "Enabled",
						Name:            "enabled",
						Type:            "bool",
						Optional:        true,
						DefaultValueStr: "true",
					},
				},
			},
		},
		{
			name: "terraform id field set multiple times",
			source: []byte(`
package testpackage

type DnsRecord struct {
	Name 			   string` + " `gen:\"name,id,required\"`" + `
	GeneratedNumber	   string` + " `gen:\"internal_id,computed\"`" + `
	Enabled 		   bool` + " `gen:\"enabled,id,optional\"`" + `
	ExplicitlyOmitted  bool` + " `gen:\"-,omit\"`" + `
}
			`),

			expectedError: errors.New(""),
		},
		{
			name: "mikrotik id field set multiple times",
			source: []byte(`
package testpackage

type DnsRecord struct {
	ID 				   string` + " `gen:\"name,id,mikrotikID,required\"`" + `
	Name 			   string` + " `gen:\"name,mikrotikID,required\"`" + `
	GeneratedNumber	   string` + " `gen:\"internal_id,computed\"`" + `
	Enabled 		   bool` + " `gen:\"enabled,optional\"`" + `
	ExplicitlyOmitted  bool` + " `gen:\"-,omit\"`" + `
}
			`),

			expectedError: errors.New(""),
		},
		{
			name: "delete id field set multiple times",
			source: []byte(`
package testpackage

type DnsRecord struct {
	ID 				   string` + " `gen:\"name,id,required\"`" + `
	Name 			   string` + " `gen:\"name,mikrotikID,deleteID,required\"`" + `
	GeneratedNumber	   string` + " `gen:\"internal_id,deleteID,computed\"`" + `
	Enabled 		   bool` + " `gen:\"enabled,optional\"`" + `
	ExplicitlyOmitted  bool` + " `gen:\"-,omit\"`" + `
}
			`),

			expectedError: errors.New(""),
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
			// todo(maksym): this condition does not check the error type since we don't have specific errors yet
			if (tc.expectedError == nil) != (err == nil) {
				t.Errorf("expected error to be %v, got %v", tc.expectedError, err)
			}
			if err != nil {
				return
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
