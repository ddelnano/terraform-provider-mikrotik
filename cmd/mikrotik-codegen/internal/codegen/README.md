MikroTik code generation
========================

This tool allows generating MikroTik resources for API client and Terraform resources based on Mikrotik struct definition.

## MikroTik client resource
To generate new MikroTik resource definition, simply run
```sh
$ go run ./cmd/mikrotik-codegen mikrotik -name BridgeVlan -commandBase "/interface/bridge/vlan"
```
where

`name` - a name of MikroTik resource to generate.

`commandBase` - base path to craft commands for CRUD operations.

It is also possible to pre-fill list of fields using `-inspect-definition-file` (see [Experimental](#experimental) section)
```sh
$ go run ./cmd/mikrotik-codegen mikrotik -name BridgeVlan -commandBase "/interface/bridge/vlan" -inspect-definition-file ./inspect_vlan.txt
```

## Terraform resource
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
$ go run ./cmd/mikrotik-codegen terraform -src client/resource.go -struct MikrotikResource > mikrotik/resource_new.go
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


## Experimental

This section contains documentation for experimental and non-stable features.

### Generate Mikrotik resource using /console/inspect definition file

Modern RouterOS versions (>7.x) provide new `/console/inspect` command to query hierarchy or syntax of particular command.

For example, `/console/inspect  path=interface,list request=child` prints `parent-child` relationship of the command:
```
Columns: TYPE, NAME, NODE-TYPE
TYPE   NAME     NODE-TYPE
self   list     dir
child  add      cmd
child  comment  cmd
child  edit     cmd
child  export   cmd
child  find     cmd
child  get      cmd
child  member   dir
child  print    cmd
child  remove   cmd
child  reset    cmd
child  set      cmd
```

while `/console/inspect  path=interface,list request=syntax` gives another set of attributes:
```
Columns: TYPE, SYMBOL, SYMBOL-TYPE, NESTED, NONORM, TEXT
TYPE    SYMBOL   SYMBOL-TYPE  NESTED  NONORM  TEXT
syntax           collection        0  yes
syntax  ..       explanation       1  no      go up to interface
syntax  add      explanation       1  no      Create a new item
syntax  comment  explanation       1  no      Set comment for items
syntax  edit     explanation       1  no
syntax  export   explanation       1  no      Print or save an export script that can be used to restore configuration
syntax  find     explanation       1  no      Find items by value
syntax  get      explanation       1  no      Gets value of item's property
syntax  member   explanation       1  no
syntax  print    explanation       1  no      Print values of item properties
syntax  remove   explanation       1  no      Remove item
syntax  reset    explanation       1  no
syntax  set      explanation       1  no      Change item properties
```

Using that information, it is possible to query (even recursively) information about all menu items and sub-commands, starting from the root `/` command.

Since this feature is recent, trying to call it with our client package results in `terminal crush`, so fully integrating `/console/inspect` into codegen binary is left for the future releases.

In general, to pre-fill list of fields during code generation, one needs:
1. Machine-readable data about available fields
2. Pass this data as `-inspect-definition-file` argument.

For step #1, we'll use this command:
```
$ :put [/console/inspect  path=interface,list,add request=child as-value]
```

which produces:
```
name=add;node-type=cmd;type=self;name=comment;node-type=arg;type=child;name=copy-from;node-type=arg;type=child;name=exclude;node-type=arg;type=child;name=include;node-type=arg;type=child;name=name;node-type=arg;type=child
```

Note, that we used `interface,list,add` as argument to `path`. The terminal equivalent would be `/interface/list/add` (not sure why it works that way, you can check [forum topic](https://forum.mikrotik.com/viewtopic.php?t=199139#p1024410))

The reason we used `add` command and not `/interface/list` menu itself, is that we need only args (fields) of `add` command - not information about possible commands for `/interface/list`

If you have `ssh` access to the Mikrotik, the following command will be helpful to get this data:
```shell
$ ssh -o Port=XXXX -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null admin@board ":put [/console inspect as-value request=child path=interface,list,add]" > inspect_definition.txt
```

After getting the definition file, just generate Mikrotik resource as usual with extra flag:
```sh
$ go run ./cmd/mikrotik-codegen mikrotik -name InterfaceList -commandBase "/interface/list" -inspect-definition-file ./inspect_definition.txt
```
and all fields for the struct will be created.
