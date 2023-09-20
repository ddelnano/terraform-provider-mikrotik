package defaultaware

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
)

// BoolAttribute creates a wrapper for schema.BoolAttribute object and generates documentation with default value.
func BoolAttribute(wrapped schema.BoolAttribute) schema.Attribute {
	return boolWrapper{wrapped}
}

func (w boolWrapper) GetDescription() string {
	desc := w.BoolAttribute.GetDescription()
	if w.Default == nil {
		return desc
	}

	resp := defaults.BoolResponse{}
	w.Default.DefaultBool(context.TODO(), defaults.BoolRequest{}, &resp)
	defaultValue := resp.PlanValue.ValueBool()
	desc = fmt.Sprintf("%s Default: `%t`.", desc, defaultValue)

	return desc
}

type boolWrapper struct {
	schema.BoolAttribute
}

var _ schema.Attribute = &boolWrapper{}
