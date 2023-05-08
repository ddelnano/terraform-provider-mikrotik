package codegen

const (
	typeString  = "String"
	typeInt64   = "Int64"
	typeList    = "List"
	typeSet     = "Set"
	typeBool    = "Bool"
	typeUnknown = "unknown"
)

type (
	basetype struct {
		typeName string
	}

	// Type represents Terraform field type to use for particular MikroTik field.
	Type interface {
		// Type returns a type name as string.
		// It must be stable for the same type.
		Name() string

		// Is checks whether two types are the same.
		Is(Type) bool
	}
)

var (
	StringType  Type = basetype{typeName: typeString}
	Int64Type   Type = basetype{typeName: typeInt64}
	ListType    Type = basetype{typeName: typeList}
	SetType     Type = basetype{typeName: typeSet}
	BoolType    Type = basetype{typeName: typeBool}
	UnknownType Type = basetype{typeName: typeUnknown}
)

func (b basetype) Name() string {
	return b.typeName
}

func (b basetype) Is(t Type) bool {
	return b.typeName == t.Name()
}
