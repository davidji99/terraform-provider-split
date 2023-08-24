package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/pmcjury/terraform-provider-split/split"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{ProviderFunc: split.New})
}
