package hiera5

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceHiera5Bool() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceHiera5BoolRead,

		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"default": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func dataSourceHiera5BoolRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Reading hiera value")

	keyName := d.Get("key").(string)
	defaultValue, defaultSet := d.GetOkExists("default")
	hiera := meta.(hiera5)

	log.Printf("[ERROR] ###################  %s  %v  %v", keyName, defaultSet, defaultValue)

	v, err := hiera.bool(keyName)
	if err != nil && !defaultSet {
		log.Printf("[DEBUG] Error reading hiera value %v", err)
		return err
	}

	d.SetId(keyName)
	if err != nil {
		d.Set("value", defaultValue.(bool))
	} else {
		d.Set("value", v)
	}

	return nil
}
