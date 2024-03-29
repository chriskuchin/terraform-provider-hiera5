package hiera5

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &Hiera5ArrayDataSource{}

type Hiera5ArrayDataSource struct {
	client hiera5
}

type Hiera5ArrayDataSourceModel struct {
	ID      types.String `tfsdk:"id"`
	Key     types.String `tfsdk:"key"`
	Value   types.List   `tfsdk:"value"`
	Default types.List   `tfsdk:"default"`
	Scope   types.Map    `tfsdk:"scope"`
}

func NewArrayDataSource() datasource.DataSource {
	return &Hiera5ArrayDataSource{}
}

func (d *Hiera5ArrayDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "hiera5_array"
}

func (d *Hiera5ArrayDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(hiera5)
}

func (d *Hiera5ArrayDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":  idAttribute,
			"key": keyAttribute,
			"value": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: valueDescription,
			},
			"default": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: defaultDescription,
			},
			"scope": scopeOverrideAttribute,
		},
	}
}

func (d *Hiera5ArrayDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data Hiera5ArrayDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	scopeOverride, diag := processScopeOverrideAttribute(ctx, data.Scope)

	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return
	}

	rawList, err := d.client.array(ctx, data.Key.String(), WithScopeOverride(scopeOverride))
	if err != nil && data.Default.IsNull() {
		resp.Diagnostics.AddAttributeError(path.Root("key"),
			"key not in data",
			"When key is unavailable and a default value is not set an error is raised")
	}

	if resp.Diagnostics.HasError() {
		return
	}

	data.ID = data.Key
	if err != nil {
		data.Value = data.Default
	} else {
		listValue := []attr.Value{}
		for _, v := range rawList {
			listValue = append(listValue, types.StringValue(v.(string)))
		}

		result, diag := types.ListValue(types.StringType, listValue)
		resp.Diagnostics.Append(diag...)
		if resp.Diagnostics.HasError() {
			return
		}

		data.Value = result
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
