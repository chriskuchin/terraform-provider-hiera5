package hiera5

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ datasource.DataSource = &Hiera5ArrayDataSource{}

type Hiera5ArrayDataSource struct {
	config Hiera5ProviderModel
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

	d.config = req.ProviderData.(Hiera5ProviderModel)
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

func (ha *Hiera5ArrayDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data Hiera5ArrayDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	log.Printf("[INFO] Reading hiera array")

	// var defaultList types.List
	// if !data.Default.IsNull() {
	// 	defaultList = data.Default
	// }
	// keyName := d.Get("key").(string)
	// rawList, defaultIsSet := d.GetOk("default")
	// hiera := meta.(hiera5)

	// v, err := hiera.array(keyName)
	// if err != nil && !defaultIsSet {
	// 	log.Printf("[DEBUG] Error reading hiera array %s", err)
	// 	return err
	// }

	data.ID = data.Key
	val, _ := types.ListValue(basetypes.StringType{}, []attr.Value{basetypes.NewStringValue("test")})
	data.Value = val

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
