package hiera5

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceHiera5Hash() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceHiera5HashRead,

		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func dataSourceHiera5HashRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Reading hiera hash")

	keyName := d.Get("key").(string)
	hiera := meta.(hiera5)

	v, err := hiera.hash(keyName)
	if err != nil {
		log.Printf("[DEBUG] Error reading hiera hash %s", err)
		return err
	}

	d.SetId(keyName)
	_ = d.Set("value", v)

	return nil
}
