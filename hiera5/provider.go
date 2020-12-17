package hiera5

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Provider is the top level function for terraform's provider API
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"config": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "hiera.yaml",
			},
			"scope": {
				Type:     schema.TypeMap,
				Default:  map[string]interface{}{},
				Optional: true,
			},
			"merge": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "first",
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"hiera5":       dataSourceHiera5(),
			"hiera5_array": dataSourceHiera5Array(),
			"hiera5_hash":  dataSourceHiera5Hash(),
			"hiera5_json":  dataSourceHiera5Json(),
			"hiera5_bool":  dataSourceHiera5Bool(),
		},

		ConfigureFunc: providerConfigure,
	}
}
func providerConfigure(data *schema.ResourceData) (interface{}, error) {
	return newHiera5(
		data.Get("config").(string),
		data.Get("scope").(map[string]interface{}),
		data.Get("merge").(string),
	), nil
}
