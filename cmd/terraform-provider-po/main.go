package main

import (
	"context"
	"flag"

	"github.com/feniix/terraform-provider-po/po"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	debugFlag := flag.Bool("debug", false, "Start provider in stand-alone debug mode.")
	flag.Parse()

	serveOpts := &plugin.ServeOpts{
		ProviderFunc: po.Provider,
	}
	if debugFlag != nil && *debugFlag {
		plugin.Debug(context.Background(), "registry.terraform.io/feniix/po", serveOpts)
	} else {
		plugin.Serve(serveOpts)
	}
}
