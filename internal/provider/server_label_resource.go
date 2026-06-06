package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/IbiliAze/terraform-provider-velia/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &ServerLabelResource{}
	_ resource.ResourceWithConfigure   = &ServerLabelResource{}
	_ resource.ResourceWithImportState = &ServerLabelResource{}
)

type ServerLabelResource struct {
	client *client.Client
}

type ServerLabelResourceModel struct {
	ServerID types.Int64  `tfsdk:"server_id"`
	Label    types.String `tfsdk:"label"`
}

func NewServerLabelResource() resource.Resource {
	return &ServerLabelResource{}
}

func (r *ServerLabelResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_label"
}

func (r *ServerLabelResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resourceSchema.Schema{
		Description: "Manages the label of a Velia server. The server itself is provisioned externally.",
		Attributes: map[string]resourceSchema.Attribute{
			"server_id": resourceSchema.Int64Attribute{
				Required:    true,
				Description: "ID of the server to label.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"label": resourceSchema.StringAttribute{
				Required:    true,
				Description: "Label to assign to the server. Set to an empty string to clear the label.",
			},
		},
	}
}

func (r *ServerLabelResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ServerLabelResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ServerLabelResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UpdateServerLabel(ctx, plan.ServerID.ValueInt64(), plan.Label.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error setting server label", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ServerLabelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ServerLabelResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	out, err := r.client.ReadServer(ctx, state.ServerID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Error reading server", err.Error())
		return
	}

	newState := ServerLabelResourceModel{
		ServerID: state.ServerID,
		Label:    types.StringValue(out.Label),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *ServerLabelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ServerLabelResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UpdateServerLabel(ctx, plan.ServerID.ValueInt64(), plan.Label.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error updating server label", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ServerLabelResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ServerLabelResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, err := r.client.UpdateServerLabel(ctx, state.ServerID.ValueInt64(), ""); err != nil {
		resp.Diagnostics.AddError("Error clearing server label", err.Error())
	}
}

func (r *ServerLabelResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	serverID, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Expected numeric server ID, got %q: %s", req.ID, err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(
		resp.State.SetAttribute(ctx, path.Root("server_id"), serverID)...,
	)
}
