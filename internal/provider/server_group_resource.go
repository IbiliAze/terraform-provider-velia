package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/IbiliAze/terraform-provider-velia/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	stringplanmodifier "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &ServerGroupResource{}
	_ resource.ResourceWithConfigure   = &ServerGroupResource{}
	_ resource.ResourceWithImportState = &ServerGroupResource{}
)

type ServerGroupResource struct {
	client *client.Client
}

type ServerGroupResourceModel struct {
	ID      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	Color   types.String `tfsdk:"color"`
	Servers types.List   `tfsdk:"servers"`
}

func NewServerGroupResource() resource.Resource {
	return &ServerGroupResource{}
}

func (r *ServerGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_group"
}

func (r *ServerGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resourceSchema.Schema{
		Description: "Creates a server group in Velia.",
		Attributes: map[string]resourceSchema.Attribute{
			"id": resourceSchema.StringAttribute{
				Computed:    true,
				Description: "Server group ID returned by Velia.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": resourceSchema.StringAttribute{
				Required:    true,
				Description: "Server group name.",
			},
			"color": resourceSchema.StringAttribute{
				Required:    true,
				Description: "Server group color.",
			},
			"servers": resourceSchema.ListAttribute{
				Required:    true,
				ElementType: types.StringType,
				Description: "List of servers.",
			},
		},
	}
}

func (r *ServerGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ServerGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ServerGroupResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var servers []string
	diags := plan.Servers.ElementsAs(ctx, &servers, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	out, err := r.client.CreateServerGroup(ctx, client.CreateServerGroupRequest{
		Name:    plan.Name.ValueString(),
		Color:   plan.Color.ValueString(),
		Servers: servers,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating server group", err.Error())
		return
	}

	serversList, diags := types.ListValueFrom(ctx, types.StringType, out.Servers)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state := ServerGroupResourceModel{
		ID:      types.StringValue(strconv.FormatInt(out.ID, 10)),
		Name:    types.StringValue(out.Name),
		Color:   types.StringValue(out.Color),
		Servers: serversList,
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ServerGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ServerGroupResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.ParseInt(state.ID.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid group ID",
			fmt.Sprintf("Expected numeric group ID, got %q: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	out, err := r.client.ReadServerGroup(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError("Error reading server group", err.Error())
		return
	}

	serversList, diags := types.ListValueFrom(ctx, types.StringType, out.Servers)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	newState := ServerGroupResourceModel{
		ID:      types.StringValue(strconv.FormatInt(out.ID, 10)),
		Name:    types.StringValue(out.Name),
		Color:   types.StringValue(out.Color),
		Servers: serversList,
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *ServerGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ServerGroupResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ServerGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ServerGroupResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.ParseInt(state.ID.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid group ID",
			fmt.Sprintf("Expected numeric group ID, got %q: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	_, err = r.client.DeleteServerGroup(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting server group", err.Error())
		return
	}
}

func (r *ServerGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if _, err := strconv.ParseInt(req.ID, 10, 64); err != nil {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Expected numeric group ID, got %q: %s", req.ID, err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(
		resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...,
	)
}
