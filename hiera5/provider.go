package hiera5

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type Hiera5Provider struct{}

type Hiera5ProviderModel struct {
	Config types.String      `tfsdk:"config"`
	Scope  map[string]string `tfsdk:"scope"`
	Merge  types.String      `tfsdk:"merge"`
}

func New() provider.Provider {
	return &Hiera5Provider{}
}

func (h *Hiera5Provider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Hiera5 Provider",
		Attributes: map[string]schema.Attribute{
			"config": schema.StringAttribute{
				Description: "The location of the hiera config file",
				Optional:    true,
			},
			"scope": schema.MapAttribute{
				ElementType: types.StringType,
				Description: "The fact variables for determining which files to merge",
				Optional:    true,
			},
			"merge": schema.StringAttribute{
				Description: "The merge strategy",
				Optional:    true,
			},
		},
	}
}

func (h *Hiera5Provider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "hiera5"
}

func (h *Hiera5Provider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data Hiera5ProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Config.IsNull() {
		data.Config = types.StringValue("hiera.yml")
	}

	if data.Merge.IsNull() {
		data.Merge = types.StringValue("first")
	}

	if data.Scope == nil {
		data.Scope = map[string]string{}
	}

	scope := map[string]interface{}{}
	for k, v := range data.Scope {
		scope[k] = v
	}
	tflog.Debug(ctx, "configuring hiera5 provider with scope", scope)

	client := hiera5{
		Config: data.Config.ValueString(),
		Scope:  scope,
		Merge:  data.Merge.ValueString(),
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (h *Hiera5Provider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}

func (h *Hiera5Provider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewArrayDataSource,
	}
}

// Provider is the top level function for terraform's provider API
// func Provider() *schema.Provider {
// 	return &schema.Provider{
// 		Schema: map[string]*schema.Schema{
// 			"config": {
// 				Type:     schema.TypeString,
// 				Optional: true,
// 				Default:  "hiera.yaml",
// 			},
// 			"scope": {
// 				Type:     schema.TypeMap,
// 				Default:  map[string]interface{}{},
// 				Optional: true,
// 			},
// 			"merge": {
// 				Type:     schema.TypeString,
// 				Optional: true,
// 				Default:  "first",
// 			},
// 		},

// 		DataSourcesMap: map[string]*schema.Resource{
// 			"hiera5":       dataSourceHiera5(),
// 			"hiera5_array": dataSourceHiera5Array(),
// 			"hiera5_hash":  dataSourceHiera5Hash(),
// 			"hiera5_json":  dataSourceHiera5Json(),
// 			"hiera5_bool":  dataSourceHiera5Bool(),
// 		},

// 		ConfigureFunc: providerConfigure,
// 	}
// }
// func providerConfigure(data *schema.ResourceData) (interface{}, error) {
// 	return newHiera5(
// 		data.Get("config").(string),
// 		data.Get("scope").(map[string]interface{}),
// 		data.Get("merge").(string),
// 	), nil
// }
