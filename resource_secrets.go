package main

import (
	"fmt"
	"log"

	"github.com/Clever/unicreds"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCredstashSecret() *schema.Resource {

	return &schema.Resource{
		Create: resourceSecretCreate,
		Read:   resourceSecretRead,
		Update: resourceSecretUpdate,
		Delete: resourceSecretDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the secret",
			},
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Value of the secret",
			},
			"context": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Encryption context for the secret",
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Version of the secrets",
			},
		},
	}
}

func resourceSecretCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	name := d.Get("name").(string)
	value := d.Get("value").(string)
	version := unicreds.PaddedInt(1)

	context := unicreds.NewEncryptionContextValue()
	for k, v := range d.Get("context").(map[string]interface{}) {
		context.Set(fmt.Sprintf("%s:%v", k, v))
	}

	err := unicreds.PutSecret(&config.TableName, config.KmsKey, name, value, version, context)
	if err != nil {
		return err
	}

	d.Set("version", version)
	d.SetId(fmt.Sprintf("%s-%s-%s", name, hash(value), version))
	return resourceSecretRead(d, meta)
}

func resourceSecretRead(d *schema.ResourceData, meta interface{}) error {
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
			return err
		}
		version = v
	}

	log.Printf("[DEBUG] Getting secret for name=%q version=%q context=%+v", name, version, context)
	out, err := unicreds.GetSecret(&config.TableName, name, version, context)
	if err != nil {
		return err
	}

	d.Set("value", out.Secret)
	d.Set("version", version)
	d.SetId(fmt.Sprintf("%s-%s-%s", name, hash(out.Secret), version))

	return nil
}

func resourceSecretUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	name := d.Get("name").(string)
	value := d.Get("value").(string)

	version := d.Get("version").(int)
	context := unicreds.NewEncryptionContextValue()
	for k, v := range d.Get("context").(map[string]interface{}) {
		context.Set(fmt.Sprintf("%s:%v", k, v))
	}
	newVersion, err := unicreds.ResolveVersion(&config.TableName, name, version)
	if err != nil {
		return err
	}

	err = unicreds.PutSecret(&config.TableName, config.KmsKey, name, value, newVersion, context)
	if err != nil {
		return err
	}

	d.Set("version", newVersion)
	d.SetId(fmt.Sprintf("%s-%s-%s", name, hash(value), newVersion))
	return nil
}

func resourceSecretDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	name := d.Get("name").(string)

	err := unicreds.DeleteSecret(&config.TableName, name)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
