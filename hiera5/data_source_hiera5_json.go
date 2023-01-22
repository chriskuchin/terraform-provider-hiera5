package hiera5

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &Hiera5JSONDataSource{}

type Hiera5JSONDataSource struct {
	client hiera5
}

type Hiera5JSONDataSourceModel struct {
	ID      types.String `tfsdk:"id"`
	Key     types.String `tfsdk:"key"`
	Value   types.String `tfsdk:"value"`
	Default types.String `tfsdk:"default"`
}

func NewJSONDataSource() datasource.DataSource {
	return &Hiera5JSONDataSource{}
}

func (hb *Hiera5JSONDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "hiera5_json"
}

func (hb *Hiera5JSONDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	hb.client = req.ProviderData.(hiera5)
}

func (hb *Hiera5JSONDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"key": schema.StringAttribute{
				Required: true,
			},
			"value": schema.StringAttribute{
				Computed: true,
			},
			"default": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (hb *Hiera5JSONDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data Hiera5JSONDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	validDefault := !data.Default.IsNull() && json.Valid([]byte(data.Default.ValueString()))

	v, err := hb.client.json(ctx, data.Key.ValueString())
	if err != nil && !validDefault {
		resp.Diagnostics.AddAttributeError(path.Root("key"),
			"key not found",
			"the key was not found and no default value was provided")
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
