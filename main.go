package main

import (
	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/aviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: aviatrix.Provider,
	})
}
