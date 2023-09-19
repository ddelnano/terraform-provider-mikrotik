package defaultaware

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
)

// Int64Attribute creates a wrapper for schema.Int64Attribute object and generates documentation with default value.
func Int64Attribute(wrapped schema.Int64Attribute) schema.Attribute {
	return int64Wrapper{wrapped}
}

func (w int64Wrapper) GetDescription() string {
	desc := w.Int64Attribute.GetDescription()
	if w.Default == nil {
		return desc
	}

	resp := defaults.Int64Response{}
	w.Default.DefaultInt64(context.TODO(), defaults.Int64Request{}, &resp)
	defaultValue := resp.PlanValue.ValueInt64()
	desc = fmt.Sprintf("%s Default: `%d`.", desc, defaultValue)

	return desc
}

type int64Wrapper struct {
	schema.Int64Attribute
}

var _ schema.Attribute = &int64Wrapper{}
