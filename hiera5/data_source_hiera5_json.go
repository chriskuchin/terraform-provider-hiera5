package hiera5

import (
	"encoding/json"
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
			"default": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceHiera5JsonRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Reading hiera json")

	keyName := d.Get("key").(string)
	defaultValue := d.Get("default").(string)
	validDefault := json.Valid([]byte(defaultValue))
	hiera := meta.(hiera5)

	v, err := hiera.json(keyName)
	if err != nil && (defaultValue == "" || !validDefault) {
		log.Printf("[DEBUG] Error reading hiera json %s", err)
		return err
	}

	d.SetId(keyName)

	if err != nil && defaultValue != "" && validDefault {
		d.Set("value", defaultValue)
	} else {
		d.Set("value", v)
	}

	return nil
}
