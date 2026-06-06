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
	_ datasource.DataSource              = &NetworkDataSource{}
	_ datasource.DataSourceWithConfigure = &NetworkDataSource{}
)

type NetworkDataSource struct {
	client *client.Client
}

type NetworkDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	CIDR         types.String `tfsdk:"cidr"`
	Version      types.Int64  `tfsdk:"version"`
	Gateway      types.String `tfsdk:"gateway"`
	Network      types.String `tfsdk:"network"`
	Broadcast    types.String `tfsdk:"broadcast"`
	Netmask      types.String `tfsdk:"netmask"`
	PrefixLength types.Int64  `tfsdk:"prefix_length"`
	Registry     types.String `tfsdk:"registry"`
	City         types.String `tfsdk:"city"`
	Country      types.String `tfsdk:"country"`
	Resolver     types.String `tfsdk:"resolver"`

	// Filter attributes (optional inputs to narrow the lookup)
	FilterCIDR types.String `tfsdk:"filter_cidr"`
	FilterIP   types.String `tfsdk:"filter_ip"`
}

func NewNetworkDataSource() datasource.DataSource {
	return &NetworkDataSource{}
}

func (d *NetworkDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network"
}

func (d *NetworkDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasourceSchema.Schema{
		Description: "Reads a Velia network. Use filter_cidr or filter_ip to select a specific network.",
		Attributes: map[string]datasourceSchema.Attribute{
			"id": datasourceSchema.StringAttribute{
				Computed:    true,
				Description: "Network ID.",
			},
			"cidr": datasourceSchema.StringAttribute{
				Computed:    true,
				Description: "Network CIDR.",
			},
			"version": datasourceSchema.Int64Attribute{
				Computed:    true,
				Description: "IP version (4 or 6).",
			},
			"gateway": datasourceSchema.StringAttribute{
				Computed:    true,
				Description: "Gateway IP address.",
			},
			"network": datasourceSchema.StringAttribute{
				Computed:    true,
				Description: "Network address.",
			},
			"broadcast": datasourceSchema.StringAttribute{
				Computed:    true,
				Description: "Broadcast address.",
			},
			"netmask": datasourceSchema.StringAttribute{
				Computed:    true,
				Description: "Subnet mask.",
			},
			"prefix_length": datasourceSchema.Int64Attribute{
				Computed:    true,
				Description: "Prefix length.",
			},
			"registry": datasourceSchema.StringAttribute{
				Computed:    true,
				Description: "WHOIS registry (e.g. RIPE NCC).",
			},
			"city": datasourceSchema.StringAttribute{
				Computed:    true,
				Description: "Data center city.",
			},
			"country": datasourceSchema.StringAttribute{
				Computed:    true,
				Description: "Data center country (ISO 3166-1 code).",
			},
			"resolver": datasourceSchema.StringAttribute{
				Computed:    true,
				Description: "IP of the responsible DNS server.",
			},
			"filter_cidr": datasourceSchema.StringAttribute{
				Optional:    true,
				Description: "Filter: exact CIDR of the network to look up (e.g. 192.168.190.96/28).",
			},
			"filter_ip": datasourceSchema.StringAttribute{
				Optional:    true,
				Description: "Filter: any IP within the network to look up.",
			},
		},
	}
}

func (d *NetworkDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *NetworkDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config NetworkDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	networks, err := d.client.ListNetworks(ctx, client.ListNetworksFilter{
		CIDR: config.FilterCIDR.ValueString(),
		IP:   config.FilterIP.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error listing networks", err.Error())
		return
	}

	if len(networks) == 0 {
		resp.Diagnostics.AddError("No network found", "No networks matched the provided filters.")
		return
	}

	if len(networks) > 1 {
		resp.Diagnostics.AddError(
			"Multiple networks found",
			fmt.Sprintf("%d networks matched; use filter_cidr or filter_ip to narrow the result.", len(networks)),
		)
		return
	}

	n := networks[0]
	state := NetworkDataSourceModel{
		ID:           types.StringValue(n.ID),
		CIDR:         types.StringValue(n.CIDR),
		Version:      types.Int64Value(int64(n.Version)),
		Gateway:      types.StringValue(n.Gateway),
		Network:      types.StringValue(n.Network),
		Broadcast:    types.StringValue(n.Broadcast),
		Netmask:      types.StringValue(n.Netmask),
		PrefixLength: types.Int64Value(int64(n.PrefixLength)),
		Registry:     types.StringValue(n.Registry),
		City:         types.StringValue(n.City),
		Country:      types.StringValue(n.Country),
		Resolver:     types.StringValue(n.Resolver),
		FilterCIDR:   config.FilterCIDR,
		FilterIP:     config.FilterIP,
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
