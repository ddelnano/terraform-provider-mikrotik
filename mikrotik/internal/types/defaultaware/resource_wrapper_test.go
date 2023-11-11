package defaultaware

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestResource_wrappers(t *testing.T) {
	testCases := []struct {
		name                      string
		resourceConstructor       func() resource.Resource
		expectedWrappedAttributes map[string]bool
	}{
		{
			name: "wrapping attributes with default-aware wrappers",
			expectedWrappedAttributes: map[string]bool{
				"wrapped_string": true,
				"wrapped_bool":   true,
				"wrapped_int64":  true,
			},
			resourceConstructor: func() resource.Resource {
				return dummyResource{
					resourceSchema: schema.Schema{
						Description: "Resource description",
						Attributes: map[string]schema.Attribute{
							"wrapped_string": schema.StringAttribute{
								Optional:    true,
								Computed:    true,
								Default:     stringdefault.StaticString("*0"),
								Description: "String attribute.",
							},
							"wrapped_int64": schema.Int64Attribute{
								Optional:    true,
								Computed:    true,
								Default:     int64default.StaticInt64(42),
								Description: "Int64 attribute.",
							},
							"string_without_default": schema.StringAttribute{
								Optional:    true,
								Computed:    true,
								Description: "String attribute.",
							},
							"int64_without_default": schema.Int64Attribute{
								Optional:    true,
								Computed:    true,
								Description: "Int64 attribute.",
							},
							"wrapped_bool": schema.BoolAttribute{
								Optional:    true,
								Computed:    true,
								Default:     booldefault.StaticBool(true),
								Description: "Bool attribute.",
							},
							"bool_without_default": schema.BoolAttribute{
								Optional:    true,
								Computed:    true,
								Description: "Bool attribute.",
							},
							"unsupported_type": schema.ListAttribute{
								Optional:    true,
								Computed:    true,
								ElementType: types.StringType,
								Description: "List attribute.",
							},
						},
					},
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			wrapped := WrapResources([]func() resource.Resource{tc.resourceConstructor})[0]()
			resp := resource.SchemaResponse{}
			wrapped.Schema(context.TODO(), resource.SchemaRequest{}, &resp)
			for name, attr := range resp.Schema.Attributes {
				switch attr.(type) {
				case stringWrapper, int64Wrapper, boolWrapper:
					if _, ok := tc.expectedWrappedAttributes[name]; !ok {
						assert.Fail(t, "wrapping error", "attribute %q should have been not wrapped, but it was", name)
					}
					continue
				}
				if _, ok := tc.expectedWrappedAttributes[name]; ok {
					assert.Fail(t, "wrapping error", "attribute %q should have been wrapped, but it wasn't", name)
				}
			}
		})
	}

}

type dummyResource struct {
	resource.Resource
	resourceSchema schema.Schema
}

func (dr dummyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = dr.resourceSchema
}
