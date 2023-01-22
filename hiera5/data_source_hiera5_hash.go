package hiera5

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &Hiera5HashDataSource{}

type Hiera5HashDataSource struct {
	client hiera5
}

type Hiera5HashDataSourceModel struct {
	ID      types.String `tfsdk:"id"`
	Key     types.String `tfsdk:"key"`
	Value   types.Map    `tfsdk:"value"`
	Default types.Map    `tfsdk:"default"`
}

func NewHashDataSource() datasource.DataSource {
	return &Hiera5HashDataSource{}
}

func (hb *Hiera5HashDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "hiera5_hash"
}

func (hb *Hiera5HashDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	hb.client = req.ProviderData.(hiera5)
}

func (hb *Hiera5HashDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"key": schema.StringAttribute{
				Required: true,
			},
			"value": schema.MapAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"default": schema.MapAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (hb *Hiera5HashDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data Hiera5HashDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	// 	rawMapDefault, defaultIsSet := d.GetOk("default")

	v, err := hb.client.hash(ctx, data.Key.ValueString())
	if err != nil && data.Default.IsNull() {
		resp.Diagnostics.AddAttributeError(path.Root("key"),
			"key not found",
			"key was not found in the data and no default value was set")
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	data.ID = data.Key
	if err != nil {
		data.Value = data.Default
	} else {
		value := map[string]attr.Value{}
		for k, v := range v {
			value[k] = types.StringValue(v.(string))
		}

		actualValue, diag := types.MapValue(types.StringType, value)
		resp.Diagnostics.Append(diag...)

		data.Value = actualValue
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
