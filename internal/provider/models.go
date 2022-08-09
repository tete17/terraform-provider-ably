package ably_control

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ably App
type AblyApp struct {
	AccountID types.String `tfsdk:"account_id"`
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Status    types.String `tfsdk:"status"`
	TLSOnly   types.Bool   `tfsdk:"tls_only"`
}