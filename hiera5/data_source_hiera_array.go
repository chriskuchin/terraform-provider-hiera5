package hiera5

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceHiera5Array() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceHiera5ArrayRead,

		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
		},
	}
}

func dataSourceHiera5ArrayRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Reading hiera array")

	keyName := d.Get("key").(string)

	hiera := meta.(hiera5)
	v, err := hiera.array(keyName)
	if err != nil {
		log.Printf("[DEBUG] Error reading hiera array %s", err)
		return err
	}

	d.SetId(keyName)
	d.Set("value", v)

	return nil
}
