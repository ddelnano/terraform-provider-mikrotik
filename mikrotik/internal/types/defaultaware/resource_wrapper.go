package defaultaware

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// WrapResources wraps the list of provider's resource contructors.
//
// Later, during actual call, the resource instance is wrapped in special proxy to replace every attribute in the schema
// with proper wrapper from "defaultsaware" package.
func WrapResources(funcs []func() resource.Resource) []func() resource.Resource {
	for i, f := range funcs {
		f := f
		funcs[i] = func() resource.Resource {
			r := resourceWrapper{f()}
			return &r
		}
	}

	return funcs
}

// Schema overrides Schema functions from the wrapped resource and makes attributes default-aware.
//
// Default-aware wrappers allows generating documentation with default values, if any.
func (r resourceWrapper) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	r.Resource.Schema(ctx, req, resp)

	for name, attr := range resp.Schema.Attributes {
		switch schemaAttr := attr.(type) {
		case schema.StringAttribute:
			if schemaAttr.Default != nil {
				resp.Schema.Attributes[name] = StringAttribute(schemaAttr)
			}
		case schema.BoolAttribute:
			if schemaAttr.Default != nil {
				resp.Schema.Attributes[name] = BoolAttribute(schemaAttr)
			}
		case schema.Int64Attribute:
			if schemaAttr.Default != nil {
				resp.Schema.Attributes[name] = Int64Attribute(schemaAttr)
			}
		}
	}
}

func (r resourceWrapper) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	rwc := r.Resource.(resource.ResourceWithConfigure)
	rwc.Configure(ctx, req, resp)
}

func (r resourceWrapper) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	rwi := r.Resource.(resource.ResourceWithImportState)
	rwi.ImportState(ctx, req, resp)
}

type resourceWrapper struct {
	resource.Resource
}

var (
	_ resource.Resource                = &resourceWrapper{}
	_ resource.ResourceWithConfigure   = &resourceWrapper{}
	_ resource.ResourceWithImportState = &resourceWrapper{}
)
