package main

import (
	"github.com/spaceshuttl/terraform-provider-discord/discord"

	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return discord.Provider()
		},
	})
}
