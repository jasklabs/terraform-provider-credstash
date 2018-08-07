package credstash

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
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
				Default:     "",
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
	version := d.Get("version").(string)
	context := unicreds.NewEncryptionContextValue()
	for k, v := range d.Get("context").(map[string]interface{}) {
		context.Set(fmt.Sprintf("%s:%v", k, v))
	}

	if version == "" {
		v, err := unicreds.GetHighestVersion(&config.TableName, name)
		if err != nil {
			if err.Error() == unicreds.ErrSecretNotFound.Error() {
				log.Printf("[DEBUG] Key not found")
				if v, ok := d.GetOk("default"); ok {
					log.Printf("[DEBUG] Using default value %v", v)
					d.Set("value", v.(string))
					return nil
				}
			}
			return err
		}
		version = v
	}

	log.Printf("[DEBUG] Getting secret for name=%q version=%q context=%+v", name, version, context)
	out, err := unicreds.GetSecret(&config.TableName, name, version, context)
	if err != nil {
		if err.Error() == unicreds.ErrSecretNotFound.Error() {
			log.Printf("[DEBUG] Key not found")
			if v, ok := d.GetOk("default"); ok {
				log.Printf("[DEBUG] Using default value %v", v)
				d.Set("value", v.(string))
				return nil
			}
		}
		return err
	}

	d.Set("value", out.Secret)
	d.Set("version", version)
	d.SetId(fmt.Sprintf("%s-%s-%s", name, hash(out.Secret), version))

	return nil
}

func hash(s string) string {
	sha := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sha[:])
}
