package hiera5

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &Hiera5HashDataSource{}

type Hiera5HashDataSource struct {
	client hiera5
}

type Hiera5HashDataSourceModel struct {
	ID      types.String `tfsdk:"id"`
	Key     types.String `tfsdk:"key"`
	Value   types.List   `tfsdk:"value"`
	Default types.List   `tfsdk:"default"`
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
		Attributes: map[string]schema.Attribute{},
	}
}

func (hb *Hiera5HashDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data Hiera5HashDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// func dataSourceHiera5Hash() *schema.Resource {
// 	return &schema.Resource{
// 		Read: dataSourceHiera5HashRead,

// 		Schema: map[string]*schema.Schema{
// 			"key": {
// 				Type:     schema.TypeString,
// 				Required: true,
// 			},
// 			"value": {
// 				Type:     schema.TypeMap,
// 				Computed: true,
// 			},
// 			"default": {
// 				Type:     schema.TypeMap,
// 				Optional: true,
// 			},
// 		},
// 	}
// }

// func dataSourceHiera5HashRead(d *schema.ResourceData, meta interface{}) error {
// 	log.Printf("[INFO] Reading hiera hash")

// 	keyName := d.Get("key").(string)
// 	rawMapDefault, defaultIsSet := d.GetOk("default")
// 	hiera := meta.(hiera5)

// 	v, err := hiera.hash(keyName)
// 	if err != nil && !defaultIsSet {
// 		log.Printf("[DEBUG] Error reading hiera hash %s", err)
// 		return err
// 	}

// 	d.SetId(keyName)
// 	if err != nil && defaultIsSet {
// 		d.Set("value", rawMapDefault.(map[string]interface{}))
// 	} else {
// 		d.Set("value", v)
// 	}

// 	return nil
// }
