package credstash

import (
	"log"

	"github.com/Clever/unicreds"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceSecret() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSecretRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the secret",
			},
			"version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Version of the secrets",
			},
			"context": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Encryption context for the secret",
			},
			"value": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Value of the secret",
			},
			"default": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Default value if key does not exist",
			},
		},
	}
}

func dataSourceSecretRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	name := d.Get("name").(string)

	var version interface{}

	if _, ok := d.GetOk("version"); !ok {
		log.Printf("[DEBUG] Version for secret %s not set", name)
		v, err := unicreds.GetHighestVersion(&config.TableName, name)
		if err != nil {
			if err != unicreds.ErrSecretNotFound {
				return err
			}
		}
		version = v
	}

	if v, ok := d.GetOk("default"); ok && version == "" {
		log.Printf("[DEBUG] Using default value %v", v)
		d.Set("value", v.(string))
		d.Set("version", "default")
		d.SetId(getID(d))
		return nil
	}

	context := getContext(d)
	log.Printf("[DEBUG] Getting secret for name=%q version=%s context=%+v", name, version.(string), context)
	out, err := unicreds.GetSecret(&config.TableName, name, version.(string), context)
	if err != nil {
		return err
	}

	d.Set("value", out.Secret)
	d.Set("version", version)
	d.SetId(getID(d))

	return nil
}

func getID(d *schema.ResourceData) string {
	return d.Get("name").(string)
}
