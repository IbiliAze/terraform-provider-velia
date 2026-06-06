package provider

///////////////////////////////////////////////////////////////////////////////////MODULES
import (
	"context"

	"github.com/IbiliAze/terraform-provider-velia/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	providerSchema "github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

//////////////////////////////////////////////////////////////////////////////////////////

var (
	_ provider.Provider = &VeliaAPIProvider{}
)

type VeliaAPIProvider struct {
	version string
}

type VeliaAPIProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	APIToken types.String `tfsdk:"api_token"`
}

func New() provider.Provider {
	return &VeliaAPIProvider{
		version: "dev",
	}
}

func (p *VeliaAPIProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "velia"
	resp.Version = p.version
}

func (p *VeliaAPIProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = providerSchema.Schema{
		Description: "Terraform provider for Velia.",
		Attributes: map[string]providerSchema.Attribute{
			"endpoint": providerSchema.StringAttribute{
				Optional:    true,
				Description: "Base URL for the Velia API.",
			},
			"api_token": providerSchema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "API token for authenticating with Velia.",
			},
		},
	}
}

func (p *VeliaAPIProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config VeliaAPIProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "https://www.velia.net/api/v1"
	apiToken := ""

	if !config.Endpoint.IsNull() && !config.Endpoint.IsUnknown() {
		endpoint = config.Endpoint.ValueString()
	}

	if !config.APIToken.IsNull() && !config.APIToken.IsUnknown() {
		apiToken = config.APIToken.ValueString()
	}

	apiClient := client.New(endpoint, apiToken)

	resp.DataSourceData = apiClient
	resp.ResourceData = apiClient
}

func (p *VeliaAPIProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewCustomerContactDataSource,
		NewServerGroupDataSource,
		NewServerDataSource,
		NewNetworkDataSource,
	}
}

func (p *VeliaAPIProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewCustomerContactResource,
		NewServerGroupResource,
		NewNetworkRdnsResource,
		NewServerLabelResource,
		NewTicketResource,
	}
}
