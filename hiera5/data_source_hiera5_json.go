package hiera5

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &Hiera5JSONDataSource{}

type Hiera5JSONDataSource struct {
	client hiera5
}

type Hiera5JSONDataSourceModel struct {
	ID      types.String `tfsdk:"id"`
	Key     types.String `tfsdk:"key"`
	Value   types.List   `tfsdk:"value"`
	Default types.List   `tfsdk:"default"`
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
		Attributes: map[string]schema.Attribute{},
	}
}

func (hb *Hiera5JSONDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data Hiera5JSONDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// func dataSourceHiera5Json() *schema.Resource {
// 	return &schema.Resource{
// 		Read: dataSourceHiera5JsonRead,

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

// func dataSourceHiera5JsonRead(d *schema.ResourceData, meta interface{}) error {
// 	log.Printf("[INFO] Reading hiera json")

// 	keyName := d.Get("key").(string)
// 	defaultValue := d.Get("default").(string)
// 	validDefault := json.Valid([]byte(defaultValue))
// 	hiera := meta.(hiera5)

// 	v, err := hiera.json(keyName)
// 	if err != nil && (defaultValue == "" || !validDefault) {
// 		log.Printf("[DEBUG] Error reading hiera json %s", err)
// 		return err
// 	}

// 	d.SetId(keyName)

// 	if err != nil && defaultValue != "" && validDefault {
// 		d.Set("value", defaultValue)
// 	} else {
// 		d.Set("value", v)
// 	}

// 	return nil
// }
