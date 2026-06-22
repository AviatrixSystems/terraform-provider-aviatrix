package main

import (
	"flag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"aviatrix.com/terraform-provider-aviatrix/aviatrix"
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
