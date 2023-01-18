package hiera5

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &Hiera5StringDataSource{}

type Hiera5StringDataSource struct {
	client hiera5
}

type Hiera5StringDataSourceModel struct {
	ID      types.String `tfsdk:"id"`
	Key     types.String `tfsdk:"key"`
	Value   types.List   `tfsdk:"value"`
	Default types.List   `tfsdk:"default"`
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
		Attributes: map[string]schema.Attribute{},
	}
}

func (hb *Hiera5StringDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data Hiera5StringDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// func dataSourceHiera5() *schema.Resource {
// 	return &schema.Resource{
// 		Read: dataSourceHiera5Read,

// 		Schema: map[string]*schema.Schema{
// 			"key": {
// 				Type:     schema.TypeString,
// 				Required: true,
// 			},
// 			"value": {
// 				Type:     schema.TypeString,
// 				Computed: true,
// 			},
// 			"default": {
// 				Type:     schema.TypeString,
// 				Optional: true,
// 			},
// 		},
// 	}
// }

// func dataSourceHiera5Read(d *schema.ResourceData, meta interface{}) error {
// 	log.Printf("[INFO] Reading hiera value")

// 	keyName := d.Get("key").(string)
// 	defaultValue := d.Get("default").(string)
// 	hiera := meta.(hiera5)

// 	v, err := hiera.value(keyName)
// 	if err != nil && defaultValue == "" {
// 		log.Printf("[DEBUG] Error reading hiera value %s", err)
// 		return err
// 	}

// 	d.SetId(keyName)
// 	if err != nil {
// 		d.Set("value", defaultValue)
// 	} else {
// 		d.Set("value", v)
// 	}

// 	return nil
// }
