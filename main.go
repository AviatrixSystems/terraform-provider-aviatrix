package main

import (
	"flag"

	"github.com/AviatrixSystems/terraform-provider-aviatrix/v3/aviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{
		Debug:        debug,
		ProviderFunc: aviatrix.Provider,
		ProviderAddr: "AviatrixSystems/aviatrix",
	}

	plugin.Serve(opts)
}
