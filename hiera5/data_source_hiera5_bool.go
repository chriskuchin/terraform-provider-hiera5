package hiera5

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &Hiera5BoolDataSource{}

type Hiera5BoolDataSource struct {
	client hiera5
}

type Hiera5BoolDataSourceModel struct {
	ID      types.String `tfsdk:"id"`
	Key     types.String `tfsdk:"key"`
	Value   types.Bool   `tfsdk:"value"`
	Default types.Bool   `tfsdk:"default"`
}

func NewBoolDataSource() datasource.DataSource {
	return &Hiera5BoolDataSource{}
}

func (hb *Hiera5BoolDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "hiera5_bool"
}

func (hb *Hiera5BoolDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	hb.client = req.ProviderData.(hiera5)
}

func (hb *Hiera5BoolDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"key": schema.StringAttribute{
				Required: true,
			},
			"default": schema.BoolAttribute{
				Optional: true,
			},
			"value": schema.BoolAttribute{
				Computed: true,
			},
		},
	}
}

func (hb *Hiera5BoolDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data Hiera5BoolDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	v, err := hb.client.bool(data.Key.ValueString())
	if err != nil && data.Default.IsNull() {
		resp.Diagnostics.AddAttributeError(path.Root("key"),
			"Key not found",
			"The key was not found in the data and the default value was not set")
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	data.ID = data.Key
	if err != nil {
		data.Value = data.Default
	} else {
		data.Value = types.BoolValue(v)
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
