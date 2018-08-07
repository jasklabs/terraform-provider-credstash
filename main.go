package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/sspinc/terraform-provider-credstash/credstash"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: credstash.Provider,
	})
}
