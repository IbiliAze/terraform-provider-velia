package provider

import (
	"context"
	"fmt"

	"github.com/IbiliAze/terraform-provider-velia/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ServerDataSource{}
	_ datasource.DataSourceWithConfigure = &ServerDataSource{}
)

type ServerDataSource struct {
	client *client.Client
}

type ServerDataSourceModel struct {
	ID         types.Int64  `tfsdk:"id"`
	Label      types.String `tfsdk:"label"`
	Ticket     types.Int64  `tfsdk:"ticket"`
	Status     types.String `tfsdk:"status"`
	Issued     types.String `tfsdk:"issued"`
	Cancelled  types.String `tfsdk:"cancelled"`
	City       types.String `tfsdk:"city"`
	Country    types.String `tfsdk:"country"`
	IPMiIP     types.String `tfsdk:"ipmi_ip"`
	Online     types.Bool   `tfsdk:"online"`
	Product    types.String `tfsdk:"product"`
	Type       types.String `tfsdk:"type"`
	DeviceType types.String `tfsdk:"device_type"`
	Tags       types.List   `tfsdk:"tags"`
	Networks   types.List   `tfsdk:"networks"`
	ServerIP   types.List   `tfsdk:"server_ip"`
}

func NewServerDataSource() datasource.DataSource {
	return &ServerDataSource{}
}

func (d *ServerDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server"
}

func (d *ServerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasourceSchema.Schema{
		Description: "Reads data about a Velia server.",
		Attributes: map[string]datasourceSchema.Attribute{
			"id": datasourceSchema.Int64Attribute{
				Required:    true,
				Description: "Server ID.",
			},
			"label": datasourceSchema.StringAttribute{
				Computed:    true,
				Description: "Custom label assigned to the server.",
			},
			"ticket": datasourceSchema.Int64Attribute{
				Computed:    true,
				Description: "Order ticket number.",
			},
			"status": datasourceSchema.StringAttribute{
				Computed:    true,
				Description: "Server status: used or expiring.",
			},
			"issued": datasourceSchema.StringAttribute{
				Computed:    true,
				Description: "Date the server was issued (YYYY-MM-DD).",
			},
			"cancelled": datasourceSchema.StringAttribute{
				Computed:    true,
				Description: "Date the server was cancelled (YYYY-MM-DD), if applicable.",
			},
			"city": datasourceSchema.StringAttribute{
				Computed:    true,
				Description: "Data center city.",
			},
			"country": datasourceSchema.StringAttribute{
				Computed:    true,
				Description: "Data center country (ISO 3166-1 code).",
			},
			"ipmi_ip": datasourceSchema.StringAttribute{
				Computed:    true,
				Description: "IPMI IP address.",
			},
			"online": datasourceSchema.BoolAttribute{
				Computed:    true,
				Description: "Whether the server is currently online.",
			},
			"product": datasourceSchema.StringAttribute{
				Computed:    true,
				Description: "Product name (e.g. HP DL360).",
			},
			"type": datasourceSchema.StringAttribute{
				Computed:    true,
				Description: "Device type.",
			},
			"device_type": datasourceSchema.StringAttribute{
				Computed:    true,
				Description: "Device sub-type.",
			},
			"tags": datasourceSchema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "Tags describing hardware components.",
			},
			"networks": datasourceSchema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "CIDR networks assigned to the server.",
			},
			"server_ip": datasourceSchema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "IP addresses assigned to the server.",
			},
		},
	}
}

func (d *ServerDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ServerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config ServerDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	out, err := d.client.ReadServer(ctx, config.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Error reading server", err.Error())
		return
	}

	tags, diags := types.ListValueFrom(ctx, types.StringType, out.Tags)
	resp.Diagnostics.Append(diags...)

	networks, diags := types.ListValueFrom(ctx, types.StringType, out.Networks)
	resp.Diagnostics.Append(diags...)

	serverIP, diags := types.ListValueFrom(ctx, types.StringType, out.ServerIP)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	state := ServerDataSourceModel{
		ID:         types.Int64Value(out.ID),
		Label:      types.StringValue(out.Label),
		Ticket:     types.Int64Value(out.Ticket),
		Status:     types.StringValue(out.Status),
		Issued:     types.StringValue(out.Issued),
		Cancelled:  types.StringValue(out.Cancelled),
		City:       types.StringValue(out.City),
		Country:    types.StringValue(out.Country),
		IPMiIP:     types.StringValue(out.IPMiIP),
		Online:     types.BoolValue(out.Online),
		Product:    types.StringValue(out.Product),
		Type:       types.StringValue(out.Type),
		DeviceType: types.StringValue(out.DeviceType),
		Tags:       tags,
		Networks:   networks,
		ServerIP:   serverIP,
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
