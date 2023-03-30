package utils

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ImportStateContextUppercaseWrapper changes the ID of the resource to upper case before passing it to wrappedFunction
//
// This wrapper is useful when resource ID is MikroTik's .id.
// Due to wierd behavior, listing via MikroTik's CLI reports lowercase .id, but find request with this id via API fails
// as it expects upper case string.
//
// Usage in resource definition.
//
// SDKv2
//
//	schema.Resource{
//		Importer: &schema.ResourceImporter{
//			StateContext: utils.ImportStateContextUppercaseWrapper(schema.ImportStatePassthroughContext),
//		}
//	}
func ImportStateContextUppercaseWrapper(wrappedFunc schema.StateContextFunc) schema.StateContextFunc {
	return func(ctx context.Context, rd *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
		rd.SetId(strings.ToUpper(rd.Id()))
		return wrappedFunc(ctx, rd, i)
	}
}
