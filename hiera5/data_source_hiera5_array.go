package hiera5

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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
			"id": schema.StringAttribute{
				Computed: true,
			},
			"key": schema.StringAttribute{
				Required: true,
			},
			"value": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"default": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
		},
	}
}

func (d *Hiera5ArrayDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data Hiera5ArrayDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	tflog.Info(ctx, "reading hiera array", map[string]interface{}{"key": data.Key.String(), "default": data.Default.Elements()})

	// keyName := d.Get("key").(string)
	// rawList, defaultIsSet := d.GetOk("default")

	rawList, err := d.client.array(data.Key.String())
	if err != nil {
		resp.Diagnostics.AddAttributeWarning(path.Root("key"),
			"Failed to find the key in the data",
			"Datasource errors when an invalid key is searched for")
	}

	if data.Default.IsNull() {
		resp.Diagnostics.AddAttributeWarning(path.Root("default"),
			"Default value is not set",
			"If default value is set and key is not found it will not error")
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
			tflog.Info(ctx, fmt.Sprintf("%v", v))
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

// func dataSourceHiera5Array() *schema.Resource {
// 	return &schema.Resource{
// 		Read: dataSourceHiera5ArrayRead,

// 		Schema: map[string]*schema.Schema{
// 			"key": {
// 				Type:     schema.TypeString,
// 				Required: true,
// 			},
// 			"value": {
// 				Type: schema.TypeList,
// 				Elem: &schema.Schema{
// 					Type: schema.TypeString,
// 				},
// 				Computed: true,
// 			},
// 			"default": {
// 				Type: schema.TypeList,
// 				Elem: &schema.Schema{
// 					Type: schema.TypeString,
// 				},
// 				Optional: true,
// 			},
// 			"scope": {
// 				Type:     schema.TypeMap,
// 				Default:  map[string]interface{}{},
// 				Optional: true,
// 			},
// 			"merge": {
// 				Type:     schema.TypeString,
// 				Optional: true,
// 			},
// 		},
// 	}
// }

// func dataSourceHiera5ArrayRead(d *schema.ResourceData, meta interface{}) error {
// }
