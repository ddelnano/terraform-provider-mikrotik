MikroTik code generation
========================

This tool allows generating MikroTik resources based on client's struct definition.

## Quickstart
Just add a `codegen` tag key to struct fields:
```go
type MikrotikResource struct{
	Id             string   `mikrotik:".id"             codegen:"id,mikrotikID,deleteID"`
	Name           string   `mikrotik:"name"            codegen:"name,required,terraformID"`
	Enabled        bool     `mikrotik:"enabled"         codegen:"enabled"`
	Items          []string `mikrotik:"items"           codegen:"items,elemType=string"`
	UpdatedAt      string   `mikrotik:"updated_at"      codegen:"updated_at,computed"`
	Unused         int      `mikrotik:"unused"          codegen:"-"`
	NotImplemented int      `mikrotik:"not_implemented" codegen:"not_implemented,omit"`
	Comment        string   `mikrotik:"comment"         codegen:"comment"`
}
```

and run:
```sh
$ go run ./cmd/mikrotik-codegen -src client/resource.go -struct MikrotikResource > mikrotik/resource_new.go
```


## Supported options

|Name|Description|
|-|-|
|terraformID|Use this field during `Read` and `Import` resource|
|mikrotikID|This field is MikroTik ID field, usually `.id`|
|deleteID|Terraform resource will use this field to delete resource|
|required|Mark field as `required` in resource schema|
|optional|Mark field as `optional` in resource schema|
|computed|Mark field as `computed` in resource schema|
|elemType|Explicitly set element type for `List` or `Set` attributes. Usage `elemType=int`|
|omit|Skip this field from code generation process|
