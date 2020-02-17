package hiera5

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceHiera5Json() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceHiera5JsonRead,

		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceHiera5JsonRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Reading hiera json")

	keyName := d.Get("key").(string)
	hiera := meta.(hiera5)

	v, err := hiera.json(keyName)
	if err != nil {
		log.Printf("[DEBUG] Error reading hiera json %s", err)
		return err
	}

	d.SetId(keyName)
	_ = d.Set("value", v)

	return nil
}
