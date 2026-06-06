package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/IbiliAze/terraform-provider-velia/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	stringplanmodifier "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &TicketResource{}
	_ resource.ResourceWithConfigure   = &TicketResource{}
	_ resource.ResourceWithImportState = &TicketResource{}
)

type TicketResource struct {
	client *client.Client
}

type TicketResourceModel struct {
	ID      types.String `tfsdk:"id"`
	Topic   types.String `tfsdk:"topic"`
	Subject types.String `tfsdk:"subject"`
	Message types.String `tfsdk:"message"`
	Servers types.List   `tfsdk:"servers"`
	Queue   types.String `tfsdk:"queue"`
	Status  types.String `tfsdk:"status"`
}

func NewTicketResource() resource.Resource {
	return &TicketResource{}
}

func (r *TicketResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ticket"
}

func (r *TicketResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resourceSchema.Schema{
		Description: "Opens a support ticket with Velia. Tickets cannot be deleted via the API; destroying this resource only removes it from Terraform state.",
		Attributes: map[string]resourceSchema.Attribute{
			"id": resourceSchema.StringAttribute{
				Computed:    true,
				Description: "Ticket ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"topic": resourceSchema.StringAttribute{
				Required:    true,
				Description: "Topic queue: velianet-billing, velianet-sales, velianet-reinstall, velianet-setup, or velianet-support.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"subject": resourceSchema.StringAttribute{
				Required:    true,
				Description: "Brief summary of the issue (max 200 characters).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"message": resourceSchema.StringAttribute{
				Required:    true,
				Description: "Full description of the issue.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"servers": resourceSchema.ListAttribute{
				Optional:    true,
				ElementType: types.Int64Type,
				Description: "Server IDs to associate with the ticket.",
				PlanModifiers: []planmodifier.List{},
			},
			"queue": resourceSchema.StringAttribute{
				Computed:    true,
				Description: "Queue the ticket was routed to.",
			},
			"status": resourceSchema.StringAttribute{
				Computed:    true,
				Description: "Current ticket status.",
			},
		},
	}
}

func (r *TicketResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TicketResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan TicketResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var servers []int64
	if !plan.Servers.IsNull() && !plan.Servers.IsUnknown() {
		resp.Diagnostics.Append(plan.Servers.ElementsAs(ctx, &servers, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	out, err := r.client.CreateTicket(ctx, client.CreateTicketRequest{
		Topic:   plan.Topic.ValueString(),
		Subject: plan.Subject.ValueString(),
		Message: plan.Message.ValueString(),
		Servers: servers,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating ticket", err.Error())
		return
	}

	state, diags := ticketToModel(ctx, out, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *TicketResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TicketResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.ParseInt(state.ID.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid ticket ID",
			fmt.Sprintf("Expected numeric ticket ID, got %q: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	out, err := r.client.ReadTicket(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError("Error reading ticket", err.Error())
		return
	}

	newState, diags := ticketToModel(ctx, out, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

// Update is never called because all mutable attributes use RequiresReplace.
func (r *TicketResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan TicketResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete removes the ticket from state. The Velia API has no delete endpoint for tickets.
func (r *TicketResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *TicketResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if _, err := strconv.ParseInt(req.ID, 10, 64); err != nil {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Expected numeric ticket ID, got %q: %s", req.ID, err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(
		resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...,
	)
}

func ticketToModel(ctx context.Context, t *client.Ticket, base TicketResourceModel) (TicketResourceModel, diag.Diagnostics) {
	serversList, diags := types.ListValueFrom(ctx, types.Int64Type, t.Servers)
	return TicketResourceModel{
		ID:      types.StringValue(strconv.FormatInt(t.ID, 10)),
		Topic:   base.Topic,
		Subject: types.StringValue(t.Subject),
		Message: base.Message,
		Servers: serversList,
		Queue:   types.StringValue(t.Queue),
		Status:  types.StringValue(t.Status),
	}, diags
}
