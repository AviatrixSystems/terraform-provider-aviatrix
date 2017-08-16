package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/terraform-providers/terraform-provider-aviatrix"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: aviatrix.Provider,
	})
}
