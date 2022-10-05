package ably_control

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	tfsdk_provider "github.com/hashicorp/terraform-plugin-framework/provider"
	tfsdk_resource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type resourceRuleAmqpType struct{}

// Get Rule Resource schema
func (r resourceRuleAmqpType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return GetRuleSchema(
		map[string]tfsdk.Attribute{
			"queue_id": {
				Type:        types.StringType,
				Required:    true,
				Description: "The ID of your Ably queue",
			},
			"headers":   GetHeaderSchema(),
			"enveloped": GetEnvelopedchema(),
			"format":    GetFormatSchema(),
		},
		"The `ably_rule_amqp` resource allows you to create and manage an Ably integration rule for AMQP. Read more at https://ably.com/docs/general/firehose/amqp-rule"), nil
}

// New resource instance
func (r resourceRuleAmqpType) NewResource(_ context.Context, p tfsdk_provider.Provider) (tfsdk_resource.Resource, diag.Diagnostics) {
	return resourceRuleAmqp{
		p: *(p.(*provider)),
	}, nil
}

type resourceRuleAmqp struct {
	p provider
}

func (r *resourceRuleAmqp) Provider() *provider {
	return &r.p
}

func (r *resourceRuleAmqp) Name() string {
	return "AMQP"
}

// Create a new resource
func (r resourceRuleAmqp) Create(ctx context.Context, req tfsdk_resource.CreateRequest, resp *tfsdk_resource.CreateResponse) {
	CreateRule[AblyRuleTargetAmqp](&r, ctx, req, resp)
}

// Read resource
func (r resourceRuleAmqp) Read(ctx context.Context, req tfsdk_resource.ReadRequest, resp *tfsdk_resource.ReadResponse) {
	ReadRule[AblyRuleTargetAmqp](&r, ctx, req, resp)
}

// // Update resource
func (r resourceRuleAmqp) Update(ctx context.Context, req tfsdk_resource.UpdateRequest, resp *tfsdk_resource.UpdateResponse) {
	UpdateRule[AblyRuleTargetAmqp](&r, ctx, req, resp)
}

// Delete resource
func (r resourceRuleAmqp) Delete(ctx context.Context, req tfsdk_resource.DeleteRequest, resp *tfsdk_resource.DeleteResponse) {
	DeleteRule[AblyRuleTargetAmqp](&r, ctx, req, resp)
}

// Import resource
func (r resourceRuleAmqp) ImportState(ctx context.Context, req tfsdk_resource.ImportStateRequest, resp *tfsdk_resource.ImportStateResponse) {
	ImportRule(&r, ctx, req, resp)
}
