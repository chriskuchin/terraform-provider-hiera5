package hiera5

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Provider is the top level funtion for terraform's provider API
func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{},
	}
}
