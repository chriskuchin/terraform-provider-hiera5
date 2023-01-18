package hiera5

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &Hiera5BoolDataSource{}

type Hiera5BoolDataSource struct {
	client hiera5
}

type Hiera5BoolDataSourceModel struct {
	ID      types.String `tfsdk:"id"`
	Key     types.String `tfsdk:"key"`
	Value   types.List   `tfsdk:"value"`
	Default types.List   `tfsdk:"default"`
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
		Attributes: map[string]schema.Attribute{},
	}
}

func (hb *Hiera5BoolDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data Hiera5BoolDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// func dataSourceHiera5Bool() *schema.Resource {
// 	return &schema.Resource{
// 		Read: dataSourceHiera5BoolRead,

// 		Schema: map[string]*schema.Schema{
// 			"key": {
// 				Type:     schema.TypeString,
// 				Required: true,
// 			},
// 			"value": {
// 				Type:     schema.TypeBool,
// 				Computed: true,
// 			},
// 			"default": {
// 				Type:     schema.TypeBool,
// 				Optional: true,
// 			},
// 		},
// 	}
// }

// func dataSourceHiera5BoolRead(d *schema.ResourceData, meta interface{}) error {
// 	log.Printf("[INFO] Reading hiera value")

// 	keyName := d.Get("key").(string)
// 	defaultValue, defaultSet := d.GetOkExists("default")
// 	hiera := meta.(hiera5)

// 	log.Printf("[INFO] ###################  %s  %v  %v", keyName, defaultSet, defaultValue)

// 	v, err := hiera.bool(keyName)
// 	if err != nil && !defaultSet {
// 		log.Printf("[ERROR] Error reading hiera value %v", err)
// 		return err
// 	}

// 	d.SetId(keyName)
// 	if err != nil {
// 		d.Set("value", defaultValue.(bool))
// 	} else {
// 		d.Set("value", v)
// 	}

// 	return nil
// }
