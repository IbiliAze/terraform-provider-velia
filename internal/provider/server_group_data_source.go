package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/IbiliAze/terraform-provider-velia/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ServerGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &ServerGroupDataSource{}
)

type ServerGroupDataSource struct {
	client *client.Client
}

type ServerGroupDataSourceModel struct {
	ID      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	Color   types.String `tfsdk:"color"`
	Servers types.List   `tfsdk:"servers"`
}

func NewServerGroupDataSource() datasource.DataSource {
	return &ServerGroupDataSource{}
}

func (d *ServerGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_group"
}

func (d *ServerGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasourceSchema.Schema{
		Description: "Reads a server group in Velia.",
		Attributes: map[string]datasourceSchema.Attribute{
			"id": datasourceSchema.StringAttribute{
				Required:    true,
				Description: "Server group ID returned by Velia.",
			},
			"name": datasourceSchema.StringAttribute{
				Computed:    true,
				Description: "Server group name.",
			},
			"color": datasourceSchema.StringAttribute{
				Computed:    true,
				Description: "Server group color.",
			},
			"servers": datasourceSchema.ListAttribute{
				Computed:    true,
				ElementType: types.Int64Type,
				Description: "List of server IDs in the group.",
			},
		},
	}
}

func (d *ServerGroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ServerGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config ServerGroupDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.ParseInt(config.ID.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid group ID",
			fmt.Sprintf("Expected numeric group ID, got %q: %s", config.ID.ValueString(), err.Error()),
		)
		return
	}

	out, err := d.client.ReadServerGroup(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError("Error reading server group", err.Error())
		return
	}

	serversList, diags := types.ListValueFrom(ctx, types.Int64Type, out.Servers)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state := ServerGroupDataSourceModel{
		ID:      types.StringValue(strconv.FormatInt(out.ID, 10)),
		Name:    types.StringValue(out.Name),
		Color:   types.StringValue(out.Color),
		Servers: serversList,
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
