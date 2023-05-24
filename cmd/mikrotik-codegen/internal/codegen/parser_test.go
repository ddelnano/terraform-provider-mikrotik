package codegen

import (
	"errors"
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
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
			name: "terraform and mikrotik id fields are parsed",
			source: []byte(`
package testpackage

type DnsRecord struct {
	ID	 			   string` + " `codegen:\"id,mikrotikID\"`" + `
	Name 			   string` + " `codegen:\"name,required,terraformID\"`" + `
	GeneratedNumber	   string` + " `codegen:\"internal_id,computed\"`" + `
	Enabled 		   bool` + " `codegen:\"enabled,optional\"`" + `
	Omitted			   bool` + " `codegen:\"-\"`" + `
	ExplicitlyOmitted  bool` + " `codegen:\"-,omit\"`" + `
}
			`),

			expected: &Struct{
				Name:             "DnsRecord",
				TerraformIDField: "Name",
				MikrotikIDField:  "ID",
				Fields: []*Field{
					{
						OriginalName: "ID",
						Name:         "id",
						Type:         "string",
						Computed:     true,
					},
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
	ID	 			   string` + " `codegen:\"id,mikrotikID\"`" + `
	Name 			   string` + " `codegen:\"name,terraformID,deleteID,required\"`" + `
	GeneratedNumber	   string` + " `codegen:\"internal_id,computed\"`" + `
	Enabled 		   bool` + " `codegen:\"enabled,optional\"`" + `
	ExplicitlyOmitted  bool` + " `codegen:\"-\"`" + `
}
			`),

			expected: &Struct{
				Name:             "DnsRecord",
				TerraformIDField: "Name",
				MikrotikIDField:  "ID",
				DeleteField:      "Name",
				Fields: []*Field{
					{
						OriginalName: "ID",
						Name:         "id",
						Type:         "string",
						Computed:     true,
					},
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
			name: "mikrotikID is not set",
			source: []byte(`
package testpackage

type DnsRecord struct {
	Id 			   	   string` + " `codegen:\"id\"`" + `
	Name 			   string` + " `codegen:\"name,terraformID,required\"`" + `
	GeneratedNumber	   string` + " `codegen:\"internal_id,computed\"`" + `
	Enabled 		   bool` + " `codegen:\"enabled,id,optional\"`" + `
	ExplicitlyOmitted  bool` + " `codegen:\"-,omit\"`" + `
}
			`),

			expectedError: errors.New(""),
		},
		{
			name: "terraform id field set multiple times",
			source: []byte(`
package testpackage

type DnsRecord struct {
	Id 			   	   string` + " `codegen:\"id,mikrotikID\"`" + `
	Name 			   string` + " `codegen:\"name,terraformID,required\"`" + `
	GeneratedNumber	   string` + " `codegen:\"internal_id,computed\"`" + `
	Enabled 		   bool` + " `codegen:\"enabled,terraformID,optional\"`" + `
	ExplicitlyOmitted  bool` + " `codegen:\"-,omit\"`" + `
}
			`),

			expectedError: errors.New(""),
		},
		{
			name: "mikrotik id field set multiple times",
			source: []byte(`
package testpackage

type DnsRecord struct {
	ID 				   string` + " `codegen:\"id,mikrotikID\"`" + `
	Name 			   string` + " `codegen:\"name,mikrotikID,required\"`" + `
	GeneratedNumber	   string` + " `codegen:\"internal_id,computed\"`" + `
	Enabled 		   bool` + " `codegen:\"enabled,optional\"`" + `
	ExplicitlyOmitted  bool` + " `codegen:\"-,omit\"`" + `
}
			`),

			expectedError: errors.New(""),
		},
		{
			name: "delete id field set multiple times",
			source: []byte(`
package testpackage

type DnsRecord struct {
	ID 				   string` + " `codegen:\"id,required\"`" + `
	Name 			   string` + " `codegen:\"name,mikrotikID,deleteID,required\"`" + `
	GeneratedNumber	   string` + " `codegen:\"internal_id,deleteID,computed\"`" + `
	Enabled 		   bool` + " `codegen:\"enabled,optional\"`" + `
	ExplicitlyOmitted  bool` + " `codegen:\"-,omit\"`" + `
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
			assert.Equal(t, tc.expected, result)
		})
	}
}
