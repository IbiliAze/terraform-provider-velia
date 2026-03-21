package provider

///////////////////////////////////////////////////////////////////////////////////MODULES
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

//////////////////////////////////////////////////////////////////////////////////////////

var (
	_ resource.Resource                = &CustomerContactResource{}
	_ resource.ResourceWithConfigure   = &CustomerContactResource{}
	_ resource.ResourceWithImportState = &CustomerContactResource{}
)

type CustomerContactResource struct {
	client *client.Client
}

type CustomerContactResourceModel struct {
	ID    types.String `tfsdk:"id"`
	Email types.String `tfsdk:"email"`
	Type  types.String `tfsdk:"type"`
}

func NewCustomerContactResource() resource.Resource {
	return &CustomerContactResource{}
}

func (r *CustomerContactResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_customer_contact"
}

func (r *CustomerContactResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resourceSchema.Schema{
		Description: "Creates a customer contact in Velia.",
		Attributes: map[string]resourceSchema.Attribute{
			"id": resourceSchema.StringAttribute{
				Computed:    true,
				Description: "Contact ID returned by Velia.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"email": resourceSchema.StringAttribute{
				Required:    true,
				Description: "Contact email address.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": resourceSchema.StringAttribute{
				Required:    true,
				Description: "Contact type, for example billing or technical.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *CustomerContactResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CustomerContactResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CustomerContactResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	out, err := r.client.CreateCustomerContact(ctx, client.CreateCustomerContactRequest{
		Email: plan.Email.ValueString(),
		Type:  plan.Type.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating customer contact", err.Error())
		return
	}

	state := CustomerContactResourceModel{
		ID:    types.StringValue(strconv.FormatInt(out.ID, 10)),
		Email: types.StringValue(out.Email),
		Type:  types.StringValue(out.Type),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *CustomerContactResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state CustomerContactResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.ParseInt(state.ID.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid contact ID",
			fmt.Sprintf("Expected numeric contact ID, got %q: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	out, err := r.client.ReadCustomerContact(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError("Error reading customer contact", err.Error())
		return
	}

	newState := CustomerContactResourceModel{
		ID:    types.StringValue(strconv.FormatInt(out.ID, 10)),
		Email: types.StringValue(out.Email),
		Type:  types.StringValue(out.Type),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *CustomerContactResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// email and type are marked RequiresReplace(), so Terraform should not
	// call Update for those changes.
	var plan CustomerContactResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *CustomerContactResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CustomerContactResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.ParseInt(state.ID.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid contact ID",
			fmt.Sprintf("Expected numeric contact ID, got %q: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	_, err = r.client.DeleteCustomerContact(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting customer contact", err.Error())
		return
	}

	// Successful delete: remove resource from state by not setting new state.
}

func (r *CustomerContactResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if _, err := strconv.ParseInt(req.ID, 10, 64); err != nil {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Expected numeric contact ID, got %q: %s", req.ID, err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(
		resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...,
	)
}
