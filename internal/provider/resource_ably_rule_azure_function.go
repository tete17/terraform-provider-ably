package ably_control

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	tfsdk_provider "github.com/hashicorp/terraform-plugin-framework/provider"
	tfsdk_resource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type resourceRuleAzureFunctionType struct{}

// Get Rule Resource schema
func (r resourceRuleAzureFunctionType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return GetRuleSchema(
		map[string]tfsdk.Attribute{
			"headers": GetHeaderSchema(),
			"azure_app_id": {
				Type:        types.StringType,
				Required:    true,
				Description: "The Microsoft Azure Application ID. You can find your Microsoft Azure Application ID",
			},
			"function_name": {
				Type:        types.StringType,
				Required:    true,
				Description: "The name of your Microsoft Azure Function",
			},
			"signing_key_id": {
				Type:        types.StringType,
				Optional:    true,
				Description: "The signing key ID for use in batch mode. Ably will optionally sign the payload using an API key ensuring your servers can validate the payload using the private API key. See the [webhook security docs](https://ably.com/docs/general/webhooks#security) for more information",
			},
			"format": {
				Type:     types.StringType,
				Optional: true,
			},
		},
		"The `ably_rule_azure_function` resource allows you to create and manage an Ably integration rule for Microsoft Azure Functions. Read more at https://ably.com/docs/general/webhooks/azure",
	), nil
}

// New resource instance
func (r resourceRuleAzureFunctionType) NewResource(_ context.Context, p tfsdk_provider.Provider) (tfsdk_resource.Resource, diag.Diagnostics) {
	return resourceRuleAzureFunction{
		p: *(p.(*provider)),
	}, nil
}

type resourceRuleAzureFunction struct {
	p provider
}

// Create a new resource
func (r resourceRuleAzureFunction) Create(ctx context.Context, req tfsdk_resource.CreateRequest, resp *tfsdk_resource.CreateResponse) {
	// Checks whether the provider and API Client are configured. If they are not, the provider responds with an error.
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply",
		)
		return
	}

	// Gets plan values
	var p AblyRuleDecoder[*AblyRuleTargetAzureFunction]
	diags := req.Plan.Get(ctx, &p)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan := p.Rule()
	plan_values := GetPlanRule(plan)

	// Creates a new Ably Rule by invoking the CreateRule function from the Client Library
	rule, err := r.p.client.CreateRule(plan.AppID.Value, &plan_values)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Resource",
			"Could not create resource, unexpected error: "+err.Error(),
		)
		return
	}

	response_values := GetRuleResponse(&rule, &plan)

	// Sets state for the new Ably App.
	diags = resp.State.Set(ctx, response_values)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource
func (r resourceRuleAzureFunction) Read(ctx context.Context, req tfsdk_resource.ReadRequest, resp *tfsdk_resource.ReadResponse) {
	// Gets the current state. If it is unable to, the provider responds with an error.
	var s AblyRuleDecoder[*AblyRuleTargetAzureFunction]
	diags := req.State.Get(ctx, &s)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	state := s.Rule()

	// Gets the Ably App ID and Ably Rule ID value for the resource
	app_id := s.AppID.Value
	rule_id := s.ID.Value

	// Get Rule data
	rule, _ := r.p.client.Rule(app_id, rule_id)

	response_values := GetRuleResponse(&rule, &state)

	// Sets state to app values.
	diags = resp.State.Set(ctx, &response_values)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
}

// // Update resource
func (r resourceRuleAzureFunction) Update(ctx context.Context, req tfsdk_resource.UpdateRequest, resp *tfsdk_resource.UpdateResponse) {
	// Gets plan values
	var p AblyRuleDecoder[*AblyRuleTargetAzureFunction]
	diags := req.Plan.Get(ctx, &p)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	var s AblyRuleDecoder[*AblyRuleTargetAzureFunction]
	diags = req.State.Get(ctx, &s)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	state := s.Rule()
	plan := p.Rule()

	rule_values := GetPlanRule(plan)

	// Gets the Ably App ID and Ably Rule ID value for the resource
	app_id := state.AppID.Value
	rule_id := state.ID.Value

	// Update Ably Rule
	rule, _ := r.p.client.UpdateRule(app_id, rule_id, &rule_values)

	response_values := GetRuleResponse(&rule, &plan)

	// Sets state to app values.
	diags = resp.State.Set(ctx, &response_values)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete resource
func (r resourceRuleAzureFunction) Delete(ctx context.Context, req tfsdk_resource.DeleteRequest, resp *tfsdk_resource.DeleteResponse) {
	// Gets the current state. If it is unable to, the provider responds with an error.
	var s AblyRuleDecoder[*AblyRuleTargetAzureFunction]
	diags := req.State.Get(ctx, &s)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	state := s.Rule()

	// Gets the Ably App ID and Ably Rule ID value for the resource
	app_id := state.AppID.Value
	rule_id := state.ID.Value

	err := r.p.client.DeleteRule(app_id, rule_id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Resource",
			"Could not delete resource, unexpected error: "+err.Error(),
		)
		return
	}

	// Remove resource from state
	resp.State.RemoveResource(ctx)
}

// // Import resource
func (r resourceRuleAzureFunction) ImportState(ctx context.Context, req tfsdk_resource.ImportStateRequest, resp *tfsdk_resource.ImportStateResponse) {
	// Save the import identifier in the id attribute
	// identifier should be in the format app_id,key_id
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: 'app_id,rule_id'. Got: %q", req.ID),
		)
		return
	}
	// Recent PR in TF Plugin Framework for paths but Hashicorp examples not updated - https://github.com/hashicorp/terraform-plugin-framework/pull/390
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("app_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[1])...)
}
