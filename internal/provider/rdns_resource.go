package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/IbiliAze/terraform-provider-velia/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	stringplanmodifier "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &NetworkRdnsResource{}
	_ resource.ResourceWithConfigure   = &NetworkRdnsResource{}
	_ resource.ResourceWithImportState = &NetworkRdnsResource{}
)

type NetworkRdnsResource struct {
	client *client.Client
}

type NetworkRdnsResourceModel struct {
	ID        types.String `tfsdk:"id"`
	NetworkID types.String `tfsdk:"network_id"`
	IP        types.String `tfsdk:"ip"`
	Type      types.String `tfsdk:"type"`
	Rdata     types.String `tfsdk:"rdata"`
	Version   types.Int64  `tfsdk:"version"`
}

func NewNetworkRdnsResource() resource.Resource {
	return &NetworkRdnsResource{}
}

func (r *NetworkRdnsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_rdns"
}

func (r *NetworkRdnsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resourceSchema.Schema{
		Description: "Manages a reverse DNS entry for a Velia network.",
		Attributes: map[string]resourceSchema.Attribute{
			"id": resourceSchema.StringAttribute{
				Computed:    true,
				Description: "rDNS entry ID returned by Velia.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"network_id": resourceSchema.StringAttribute{
				Required:    true,
				Description: "ID of the network that owns this IP address.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"ip": resourceSchema.StringAttribute{
				Required:    true,
				Description: "IP address for which to set the rDNS record.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": resourceSchema.StringAttribute{
				Required:    true,
				Description: "Record type: PTR, CNAME, or NS.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"rdata": resourceSchema.StringAttribute{
				Required:    true,
				Description: "Hostname or target of the rDNS record.",
			},
			"version": resourceSchema.Int64Attribute{
				Computed:    true,
				Description: "IP version (4 or 6).",
			},
		},
	}
}

func (r *NetworkRdnsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	apiClient, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}

	r.client = apiClient
}

func (r *NetworkRdnsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan NetworkRdnsResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	out, err := r.client.CreateRdns(ctx, client.CreateRdnsRequest{
		IP:    plan.IP.ValueString(),
		Type:  plan.Type.ValueString(),
		Rdata: plan.Rdata.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating rDNS entry", err.Error())
		return
	}

	state := rdnsToModel(out, plan.NetworkID.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *NetworkRdnsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state NetworkRdnsResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.ParseInt(state.ID.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid rDNS ID",
			fmt.Sprintf("Expected numeric rDNS ID, got %q: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	out, err := r.client.ReadRdns(ctx, state.NetworkID.ValueString(), id)
	if err != nil {
		resp.Diagnostics.AddError("Error reading rDNS entry", err.Error())
		return
	}

	newState := rdnsToModel(out, state.NetworkID.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *NetworkRdnsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan NetworkRdnsResourceModel
	var state NetworkRdnsResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.ParseInt(state.ID.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid rDNS ID",
			fmt.Sprintf("Expected numeric rDNS ID, got %q: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	out, err := r.client.UpdateRdns(ctx, id, client.UpdateRdnsRequest{
		Rdata: plan.Rdata.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error updating rDNS entry", err.Error())
		return
	}

	newState := rdnsToModel(out, state.NetworkID.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *NetworkRdnsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state NetworkRdnsResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.ParseInt(state.ID.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid rDNS ID",
			fmt.Sprintf("Expected numeric rDNS ID, got %q: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	if err := r.client.DeleteRdns(ctx, id); err != nil {
		resp.Diagnostics.AddError("Error deleting rDNS entry", err.Error())
	}
}

// ImportState accepts "network_id:rdns_id" format.
func (r *NetworkRdnsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, ":", 2)
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Expected format \"network_id:rdns_id\", got %q", req.ID),
		)
		return
	}

	if _, err := strconv.ParseInt(parts[1], 10, 64); err != nil {
		resp.Diagnostics.AddError(
			"Invalid rDNS ID in import ID",
			fmt.Sprintf("Expected numeric rDNS ID, got %q: %s", parts[1], err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("network_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}

func rdnsToModel(r *client.Rdns, networkID string) NetworkRdnsResourceModel {
	return NetworkRdnsResourceModel{
		ID:        types.StringValue(strconv.FormatInt(r.ID, 10)),
		NetworkID: types.StringValue(networkID),
		IP:        types.StringValue(r.IP),
		Type:      types.StringValue(r.Type),
		Rdata:     types.StringValue(r.Rdata),
		Version:   types.Int64Value(int64(r.Version)),
	}
}
