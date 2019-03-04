package discord

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	// The actual provider
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Discord Bot token for authenticating with Discord's API",
				DefaultFunc: schema.EnvDefaultFunc("DISCORD_TOKEN", nil),
			},
		},
	}
}
