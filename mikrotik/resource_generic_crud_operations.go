package mikrotik

import (
	"context"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/ddelnano/terraform-provider-mikrotik/mikrotik/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type (
	CreateFunc func(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse)
	ReadFunc   func(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse)
	UpdateFunc func(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse)
	DeleteFunc func(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse)
)

// GenericCreateResource creates the resource and sets the initial Terraform state.
//
// terraformModel and mikrotikModel must be passed as pointers
func GenericCreateResource(terraformModel interface{}, mikrotikModel client.Resource, client *client.Mikrotik) CreateFunc {
	return func(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

		diags := req.Plan.Get(ctx, terraformModel)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		if err := utils.TerraformModelToMikrotikStruct(terraformModel, mikrotikModel); err != nil {
			resp.Diagnostics.AddError("Cannot copy model: Terraform -> MikroTik", err.Error())
			return
		}

		created, err := client.Add(mikrotikModel)
		if err != nil {
			resp.Diagnostics.AddError("Creation failed", err.Error())
			return
		}

		if err := utils.MikrotikStructToTerraformModel(created, terraformModel); err != nil {
			resp.Diagnostics.AddError("Cannot copy model: MikroTik -> Terraform", err.Error())
			return
		}

		resp.Diagnostics.Append(resp.State.Set(ctx, terraformModel)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
}

// GenericReadResource refreshes the Terraform state with the latest data.
func GenericReadResource(terraformModel interface{}, mikrotikModel client.Resource, client *client.Mikrotik) ReadFunc {
	return func(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
		resp.Diagnostics.Append(req.State.Get(ctx, terraformModel)...)
		if resp.Diagnostics.HasError() {
			return
		}
		if err := utils.TerraformModelToMikrotikStruct(terraformModel, mikrotikModel); err != nil {
			resp.Diagnostics.AddError("Cannot copy model: Terraform -> MikroTik", err.Error())
			return
		}

		resource, err := client.Find(mikrotikModel)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading remote resource",
				err.Error(),
			)
			return
		}
		if err := utils.MikrotikStructToTerraformModel(resource, terraformModel); err != nil {
			resp.Diagnostics.AddError("Cannot copy model: MikroTik -> Terraform", err.Error())
			return
		}

		resp.Diagnostics.Append(resp.State.Set(ctx, &terraformModel)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
}

// GenericUpdateResource updates the resource and sets the updated Terraform state on success.
func GenericUpdateResource(terraformModel interface{}, mikrotikModel client.Resource, client *client.Mikrotik) UpdateFunc {
	return func(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
		resp.Diagnostics.Append(req.Plan.Get(ctx, terraformModel)...)
		if resp.Diagnostics.HasError() {
			return
		}
		if err := utils.TerraformModelToMikrotikStruct(terraformModel, mikrotikModel); err != nil {
			resp.Diagnostics.AddError("Cannot copy model: Terraform -> MikroTik", err.Error())
			return
		}
		updated, err := client.Update(mikrotikModel)
		if err != nil {
			resp.Diagnostics.AddError("Update failed", err.Error())
			return
		}
		if err := utils.MikrotikStructToTerraformModel(updated, terraformModel); err != nil {
			resp.Diagnostics.AddError("Cannot copy model: MikroTik -> Terraform", err.Error())
			return
		}

		resp.Diagnostics.Append(resp.State.Set(ctx, terraformModel)...)
	}
}

// GenericDeleteResource deletes the resource and removes the Terraform state on success.
func GenericDeleteResource(terraformModel interface{}, mikrotikModel client.Resource, client *client.Mikrotik) DeleteFunc {
	return func(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
		resp.Diagnostics.Append(req.State.Get(ctx, terraformModel)...)
		if resp.Diagnostics.HasError() {
			return
		}

		if err := utils.TerraformModelToMikrotikStruct(terraformModel, mikrotikModel); err != nil {
			resp.Diagnostics.AddError("Cannot copy model: Terraform -> MikroTik", err.Error())
			return
		}

		if err := client.Delete(mikrotikModel); err != nil {
			resp.Diagnostics.AddError("Could not delete MikroTik resource", err.Error())
			return
		}
	}
}
