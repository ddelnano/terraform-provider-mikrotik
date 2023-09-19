package defaultaware

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
)

// StringAttribute creates a wrapper for schema.StringAttribute object and generates documentation with default value.
func StringAttribute(wrapped schema.StringAttribute) schema.Attribute {
	return stringWrapper{wrapped}
}

func (w stringWrapper) GetDescription() string {
	desc := w.StringAttribute.GetDescription()
	if w.Default == nil {
		return desc
	}

	resp := defaults.StringResponse{}
	w.Default.DefaultString(context.TODO(), defaults.StringRequest{}, &resp)
	defaultValue := resp.PlanValue.ValueString()
	if defaultValue == "" {
		defaultValue = `""`
	}
	desc = fmt.Sprintf("%s Default: `%s`.", desc, defaultValue)

	return desc
}

type stringWrapper struct {
	schema.StringAttribute
}

var _ schema.Attribute = &stringWrapper{}
