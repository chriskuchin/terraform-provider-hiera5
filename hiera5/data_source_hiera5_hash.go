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
			"default": {
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
	}
}

func dataSourceHiera5HashRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Reading hiera hash")

	keyName := d.Get("key").(string)
	rawMapDefault, defaultIsSet := d.GetOk("default")
	hiera := meta.(hiera5)

	v, err := hiera.hash(keyName)
	if err != nil && !defaultIsSet {
		log.Printf("[DEBUG] Error reading hiera hash %s", err)
		return err
	}

	d.SetId(keyName)
	if err != nil && defaultIsSet {
		d.Set("value", rawMapDefault.(map[string]interface{}))
	} else {
		d.Set("value", v)
	}

	return nil
}
