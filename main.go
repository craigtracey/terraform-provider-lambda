package main

import (
    "github.com/craigtracey/terraform-provider-lambda/pkg/provider"
    "github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
    plugin.Serve(&plugin.ServeOpts{
        ProviderFunc: provider.Provider,
    })
}
