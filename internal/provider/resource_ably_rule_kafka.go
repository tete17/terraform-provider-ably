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

type resourceRuleKafkaType struct{}

// Get Rule Resource schema
func (r resourceRuleKafkaType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return GetRuleSchema(
		map[string]tfsdk.Attribute{
			"routing_key": {
				Type:        types.StringType,
				Required:    true,
				Description: "The Kafka partition key. This is used to determine which partition a message should be routed to, where a topic has been partitioned. routingKey should be in the format topic:key where topic is the topic to publish to, and key is the value to use as the message key",
			},
			"enveloped": {
				Type:        types.BoolType,
				Optional:    true,
				Description: "Delivered messages are wrapped in an Ably envelope by default that contains metadata about the message and its payload. The form of the envelope depends on whether it is part of a Webhook/Function or a Queue/Firehose rule. For everything besides Webhooks, you can ensure you only get the raw payload by unchecking `Enveloped` when setting up the rule",
			},
			"format": {
				Type:        types.StringType,
				Optional:    true,
				Description: "JSON provides a simpler text-based encoding, whereas MsgPack provides a more efficient binary encoding",
			},
			"brokers": {
				Type: types.ListType{
					ElemType: types.StringType,
				},
				Required:    true,
				Description: "This is a list of brokers that host your Kafka partitions. Each broker is specified using the format `host`, `host:port` or `ip:port`",
			},
			"auth": {
				Required:    true,
				Description: "The Kafka [authentication mechanism](https://docs.confluent.io/platform/current/kafka/overview-authentication-methods.html)",
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"sasl": {
						Optional:    true,
						Description: "SASL(Simple Authentication Security Layer) / SCRAM (Salted Challenge Response Authentication Mechanism) uses usernames and passwords stored in ZooKeeper. Credentials are created during installation. See documentation on [configuring SCRAM](https://docs.confluent.io/platform/current/kafka/authentication_sasl/authentication_sasl_scram.html#kafka-sasl-auth-scram)",
						Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
							"mechanism": {
								Description: "`plain` `scram-sha-256` `scram-sha-512`. The hash type to use. SCRAM supports either SHA-256 or SHA-512 hash functions",
								Type:        types.StringType,
								Required:    true,
							},
							"username": {
								Description: "Kafka login credential",
								Type:        types.StringType,
								Required:    true,
								Sensitive:   true,
							},
							"password": {
								Description: "Kafka login credential",
								Type:        types.StringType,
								Required:    true,
								Sensitive:   true,
							},
						}),
					},
				}),
			},
		},
		"The `ably_rule_kafka` resource allows you to create and manage an Ably integration rule for Kafka. Read more at https://ably.com/docs/general/firehose/kafka-rule",
	), nil
}

// New resource instance
func (r resourceRuleKafkaType) NewResource(_ context.Context, p tfsdk_provider.Provider) (tfsdk_resource.Resource, diag.Diagnostics) {
	return resourceRuleKafka{
		p: *(p.(*provider)),
	}, nil
}

type resourceRuleKafka struct {
	p provider
}

// Create a new resource
func (r resourceRuleKafka) Create(ctx context.Context, req tfsdk_resource.CreateRequest, resp *tfsdk_resource.CreateResponse) {
	// Checks whether the provider and API Client are configured. If they are not, the provider responds with an error.
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply",
		)
		return
	}

	// Gets plan values
	var p AblyRuleDecoder[*AblyRuleTargetKafka]
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
func (r resourceRuleKafka) Read(ctx context.Context, req tfsdk_resource.ReadRequest, resp *tfsdk_resource.ReadResponse) {
	// Gets the current state. If it is unable to, the provider responds with an error.
	var s AblyRuleDecoder[*AblyRuleTargetKafka]
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
func (r resourceRuleKafka) Update(ctx context.Context, req tfsdk_resource.UpdateRequest, resp *tfsdk_resource.UpdateResponse) {
	// Gets plan values
	var p AblyRuleDecoder[*AblyRuleTargetKafka]
	diags := req.Plan.Get(ctx, &p)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	var s AblyRuleDecoder[*AblyRuleTargetKafka]
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
func (r resourceRuleKafka) Delete(ctx context.Context, req tfsdk_resource.DeleteRequest, resp *tfsdk_resource.DeleteResponse) {
	// Gets the current state. If it is unable to, the provider responds with an error.
	var s AblyRuleDecoder[*AblyRuleTargetKafka]
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
func (r resourceRuleKafka) ImportState(ctx context.Context, req tfsdk_resource.ImportStateRequest, resp *tfsdk_resource.ImportStateResponse) {
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
