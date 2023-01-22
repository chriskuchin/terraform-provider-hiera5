package hiera5

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
				Description: "The location of the hiera config file. Default: ./hiera.yml",
				Optional:    true,
			},
			"scope": scopeAttribute,
			"merge": schema.StringAttribute{
				MarkdownDescription: "The merge strategy to use in merging data. Possible values include `first`, `unique`, `hash`, and `deep`. Further documentation can be found [here](https://www.puppet.com/docs/puppet/7/hiera_merging.html). Default: first",
				Optional:            true,
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
		NewBoolDataSource,
		NewStringDataSource,
		NewJSONDataSource,
		NewHashDataSource,
	}
}
