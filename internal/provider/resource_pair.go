// Copyright (c) Persona
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ resource.ResourceWithModifyPlan = (*PairResource)(nil)

type PlanOrState interface {
	Set(context.Context, interface{}) diag.Diagnostics
}

func NewPairResource() resource.Resource {
	return &PairResource{}
}

type PairResource struct{}

func (r *PairResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var model pairModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...)

	if resp.Diagnostics.HasError() {
		return
	}

	model.ID = types.StringValue("-")

	r.modify(ctx, model, map[string]string{}, &resp.Diagnostics, &resp.State)
}

// Delete does not need to explicitly call resp.State.RemoveResource() as this is automatically handled by the
// [framework](https://github.com/hashicorp/terraform-plugin-framework/pull/301).
func (r *PairResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *PairResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pair"
}

func (r *PairResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Will be when the resource is being deleted.
	if req.Plan.Raw.IsNull() {
		return
	}

	var model pairModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read existing result field from state, if present.
	existingResult := make(map[string]types.String)
	if !req.State.Raw.IsNull() {
		resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("result"), &existingResult)...)

		if resp.Diagnostics.HasError() {
			return
		}
	}

	convertedExistingResult := make(map[string]string, len(existingResult))
	for key, value := range existingResult {
		convertedExistingResult[key] = value.ValueString()
	}

	r.modify(ctx, model, convertedExistingResult, &resp.Diagnostics, &resp.Plan)
}

// Read does not need to perform any operations as the state in ReadResourceResponse is already populated.
func (r *PairResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

func (r *PairResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Generates a mapping of keys to values that stays stable between applies and makes minimal changes when the set of keys or values changes.",
		Attributes: map[string]schema.Attribute{
			"keys": schema.SetAttribute{
				Description: "The set of keys to assign a value. An unknown key that can be assigned a value (either known or unknown) will trigger the result to be unknown.",
				ElementType: types.StringType,
				Required:    true,
			},
			"values": schema.SetAttribute{
				Description: "The set of values to assign to keys.",
				ElementType: types.StringType,
				Required:    true,
			},

			// Computed
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "A static value used internally by Terraform, this should not be referenced in configurations.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"result": schema.MapAttribute{
				Computed:    true,
				Description: "The stable mapping of keys to values, size will be the smaller of the size of keys and values. The value will generally be known at plan time unless an unknown key can be assigned a value in which the whole result will be unknown but the end result will still be stable.",
				ElementType: types.StringType,
			},
		},
	}
}

// Update ensures the plan value is copied to the state to complete the update.
func (r *PairResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var model pairModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read existing result field from state.
	existingResult := make(map[string]types.String)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("result"), &existingResult)...)

	if resp.Diagnostics.HasError() {
		return
	}

	convertedExistingResult := make(map[string]string, len(existingResult))
	for key, value := range existingResult {
		convertedExistingResult[key] = value.ValueString()
	}

	r.modify(ctx, model, convertedExistingResult, &resp.Diagnostics, &resp.State)
}

func (r *PairResource) modify(ctx context.Context, model pairModel, existingResult map[string]string, diagnostics *diag.Diagnostics, state PlanOrState) {
	keys := make([]basetypes.StringValue, len(model.Keys.Elements()))
	diagnostics.Append(model.Keys.ElementsAs(ctx, &keys, false)...)
	if diagnostics.HasError() {
		return
	}

	values := make([]basetypes.StringValue, len(model.Values.Elements()))
	diagnostics.Append(model.Values.ElementsAs(ctx, &values, false)...)
	if diagnostics.HasError() {
		return
	}

	result := pairStable(existingResult, keys, values)

	var diags diag.Diagnostics
	model.Result, diags = types.MapValueFrom(ctx, types.StringType, result)
	diagnostics.Append(diags...)
	if diagnostics.HasError() {
		return
	}

	diagnostics.Append(state.Set(ctx, model)...)
	if diagnostics.HasError() {
		return
	}
}

type pairModel struct {
	ID     types.String `tfsdk:"id"`
	Keys   types.Set    `tfsdk:"keys"`
	Result types.Map    `tfsdk:"result"`
	Values types.Set    `tfsdk:"values"`
}

func pairStable(existingResult map[string]string, keys, values []basetypes.StringValue) basetypes.MapValue {
	// First up, make a map each of keys and values to allow for easy logic below.
	keyMapping := make(map[string]bool)
	keysUnknown := 0
	valueMapping := make(map[string]bool)
	valuesUnknown := 0

	for _, key := range keys {
		if key.IsUnknown() {
			keysUnknown += 1
		} else {
			keyMapping[key.ValueString()] = true
		}
	}

	for _, value := range values {
		if value.IsUnknown() {
			valuesUnknown += 1
		} else {
			valueMapping[value.ValueString()] = true
		}
	}

	// Given an existing mapping, determine which of those should persist. If a key
	// is no longer present, no value needs to be assigned. However, if a value is
	// no longer present, a new one needs to be assigned. The latter is easily
	// achieved by leaving it out of the trimmed mapping and then allowing the
	// logic below for new keys take care of that.
	finalMapping := make(map[string]attr.Value)
	valuesUsed := make(map[string]bool)

	for key, value := range existingResult {
		if _, ok := keyMapping[key]; !ok {
			continue
		}

		if _, ok := valueMapping[value]; !ok {
			continue
		}

		finalMapping[key] = basetypes.NewStringValue(value)
		valuesUsed[value] = true
	}

	// Next, find new values for new keys (or existing ones who lost their value).
	for _, key := range keys {
		if key.IsUnknown() {
			continue
		}

		if _, ok := finalMapping[key.ValueString()]; ok {
			continue
		}

		for _, value := range values {
			if value.IsUnknown() {
				continue
			}

			if _, ok := valuesUsed[value.ValueString()]; ok {
				continue
			}

			finalMapping[key.ValueString()] = value
			valuesUsed[value.ValueString()] = true
			break
		}

		if valuesUnknown > 0 {
			if _, ok := finalMapping[key.ValueString()]; !ok {
				finalMapping[key.ValueString()] = basetypes.NewStringUnknown()
				valuesUnknown -= 1
				continue
			}
		}
	}

	// If at the end of all of this, we have some unknown keys that would map to
	// some unknown values, we sadly have to return an entirely unknown result due
	// the requirement that maps have string values.
	if keysUnknown > 0 && (valuesUnknown > 0 || (len(valueMapping)-len(valuesUsed) > 0)) {
		return basetypes.NewMapUnknown(types.StringType)
	}

	return basetypes.NewMapValueMust(types.StringType, finalMapping)
}
