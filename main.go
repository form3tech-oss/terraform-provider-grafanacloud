package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/form3tech-oss/terraform-provider-grafanacloud/grafanacloud"
)

var (
	version string = "dev"
	commit  string = ""
)

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{ProviderFunc: grafanacloud.NewProvider(version)}

	if debugMode {
		err := plugin.Debug(context.Background(), grafanacloud.Addr, opts)
		if err != nil {
			log.Fatal(err.Error())
		}
		return
	}

	plugin.Serve(opts)
}
