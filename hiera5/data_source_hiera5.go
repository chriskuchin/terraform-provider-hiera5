package hiera5

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &Hiera5StringDataSource{}

type Hiera5StringDataSource struct {
	client hiera5
}

type Hiera5StringDataSourceModel struct {
	ID      types.String `tfsdk:"id"`
	Key     types.String `tfsdk:"key"`
	Value   types.String `tfsdk:"value"`
	Default types.String `tfsdk:"default"`
}

func NewStringDataSource() datasource.DataSource {
	return &Hiera5StringDataSource{}
}

func (hb *Hiera5StringDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "hiera5"
}

func (hb *Hiera5StringDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	hb.client = req.ProviderData.(hiera5)
}

func (hb *Hiera5StringDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":  idAttribute,
			"key": keyAttribute,
			"value": schema.StringAttribute{
				Computed:    true,
				Description: valueDescription,
			},
			"default": schema.StringAttribute{
				Optional:    true,
				Description: defaultDescription,
			},
		},
	}
}

func (hb *Hiera5StringDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data Hiera5StringDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	v, err := hb.client.value(ctx, data.Key.ValueString())
	if err != nil && data.Default.IsNull() {
		resp.Diagnostics.AddAttributeError(path.Root("key"),
			"key not found",
			"the value was not found and the default value was not set")
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	data.ID = data.Key
	if err != nil {
		data.Value = data.Default
	} else {
		data.Value = types.StringValue(v)
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
