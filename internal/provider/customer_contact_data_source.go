package provider

///////////////////////////////////////////////////////////////////////////////////MODULES
import (
	"context"
	"fmt"
	"strconv"

	"github.com/IbiliAze/terraform-provider-velia/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

//////////////////////////////////////////////////////////////////////////////////////////

var (
	_ datasource.DataSource              = &CustomerContactDataSource{}
	_ datasource.DataSourceWithConfigure = &CustomerContactDataSource{}
)

type CustomerContactDataSource struct {
	client *client.Client
}

type CustomerContactDataSourceModel struct {
	ID    types.String `tfsdk:"id"`
	Email types.String `tfsdk:"email"`
	Type  types.String `tfsdk:"type"`
}

func NewCustomerContactDataSource() datasource.DataSource {
	return &CustomerContactDataSource{}
}

func (d *CustomerContactDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_customer_contact"
}

func (d *CustomerContactDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasourceSchema.Schema{
		Description: "Reads a customer contact from Velia.",
		Attributes: map[string]datasourceSchema.Attribute{
			"id": datasourceSchema.StringAttribute{
				Required:    true,
				Description: "Numeric customer contact ID.",
			},
			"email": datasourceSchema.StringAttribute{
				Computed:    true,
				Description: "Customer contact email address.",
			},
			"type": datasourceSchema.StringAttribute{
				Computed:    true,
				Description: "Customer contact type, for example billing or technical.",
			},
		},
	}
}

func (d *CustomerContactDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	apiClient, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}

	d.client = apiClient
}

func (d *CustomerContactDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config CustomerContactDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.ParseInt(config.ID.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid contact ID",
			fmt.Sprintf("Expected numeric contact ID, got %q: %s", config.ID.ValueString(), err.Error()),
		)
		return
	}

	out, err := d.client.ReadCustomerContact(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError("Error reading customer contact", err.Error())
		return
	}

	state := CustomerContactDataSourceModel{
		ID:    types.StringValue(strconv.FormatInt(out.ID, 10)),
		Email: types.StringValue(out.Email),
		Type:  types.StringValue(out.Type),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
